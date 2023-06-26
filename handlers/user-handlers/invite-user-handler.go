package user_handles

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"keycloak-user-service/types"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

func (uc *UserContext) InviteUsers() {
	var inviteUserData types.InviteUsers
	err := uc.ginCtx.Bind(&inviteUserData)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("user: cannot parse Invite User Payload: %v", inviteUserData))
		uc.ginCtx.JSON(http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(inviteUserData)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("user: cannot validate Invite User Payload: %v", inviteUserData))
		log.Error().Err(err)
		uc.ginCtx.JSON(http.StatusBadRequest, err)
		return
	}

	if inviteUserData.Emails == nil || len(inviteUserData.Emails) == 0 {
		log.Debug().Msg(fmt.Sprintf("user: must include emails of users to be invited: %v", inviteUserData))
		uc.ginCtx.JSON(http.StatusBadRequest, errors.New("must include emails of users to be invited"))
		return
	}

	orgId, err := uc.getOrgIdFromToken(uc.token)
	if err != nil {
		uc.ginCtx.JSON(http.StatusBadRequest, fmt.Sprintf("failed to retrieve %s from the access token: %s", types.ORG_ID_CLAIM_NAME, err.Error()))
		return
	}
	clientId, err := getClientIdFromToken(uc.token)
	if err != nil {
		uc.ginCtx.JSON(http.StatusBadRequest, fmt.Sprintf("failed to retrieve %s from the access token: %s", types.CLIENT_ID_CLAIM_NAME, err.Error()))
		return
	}

	response, err := uc.inviteUsers(*orgId, clientId, inviteUserData)
	if err != nil && errors.Cause(err) == types.INTERNAL_SERVER_ERROR {
		uc.ginCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	} else if err != nil && (errors.Cause(err) == types.CUSTOMER_GROUP_NOT_FOUND || errors.Cause(err) == types.MORE_THAN_ONE_CUSTOMER_GROUP) {
		uc.ginCtx.JSON(http.StatusBadRequest, err.Error())
	} else {
		uc.ginCtx.JSON(response.Code, response.Message)
		return
	}
}

func (c *UserContext) inviteUsers(orgId int, clientId *string, params types.InviteUsers) (*types.InviteUsersResponse, error) {
	destinationGroupPath, err := c.destinationGroupPath(params.IsAdmin, orgId)
	if err != nil {
		return nil, err
	}
	log.Info().Msg(fmt.Sprintf("Destination group path is %s", destinationGroupPath))

	destinationGroup, err := c.client.GetGroupByPath(c.ctx, c.token, types.KEYCLOAK_REALM, destinationGroupPath)
	if err != nil {
		return nil, err
	}

	if destinationGroup == nil {
		return nil, types.NotFound(fmt.Sprintf("cannot find group at path %s", destinationGroupPath))
	}

	response := types.NewInviteUsersResponse()
	for _, email := range params.Emails {
		err = c.createUserAndSendEmail(clientId, email, *destinationGroup.Path)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("invitation to %s failed with error: %s", email, err))
			response.AddFailure(email)
		}
	}
	return response, nil
}

func (c *UserContext) destinationGroupPath(isAdmin bool, orgId int) (string, error) {
	inviteUserError := &types.InviteUsersError{}
	params := gocloak.GetGroupsParams{
		Q:                   gocloak.StringP(fmt.Sprintf("%s:%d", types.ORG_ID_ATTRIBUTE, orgId)),
		Full:                gocloak.BoolP(true),
		BriefRepresentation: gocloak.BoolP(false),
	}

	groups, err := c.client.GetGroups(c.ctx, c.token, c.realm, params)
	if err != nil {
		return "", err
	}
	if len(groups) == 0 {
		log.Debug().Msg(fmt.Sprintf("Cannot find customer group with %s attribute set to %d", types.ORG_ID_ATTRIBUTE, orgId))
		return "", errors.Wrapf(inviteUserError.CustomerGroupNotFoundError(), "cannot find customer group with %s equal to %d", types.ORG_ID_ATTRIBUTE, orgId)
	}

	if len(groups) != 1 {
		var paths []string
		for _, group := range groups {
			paths = append(paths, *group.Path)
		}

		log.Debug().Msg(fmt.Sprintf("Found more then 1 group with %s attribute set to %d: %s", types.ORG_ID_ATTRIBUTE, orgId, strings.Join(paths, ", ")))
		return "", errors.Wrapf(inviteUserError.MoreThanOneCustomerGroup(), "cannot find a single customer group with %s equal to %d, found %d", types.ORG_ID_ATTRIBUTE, orgId, len(groups))
	}

	customerGroup := groups[0]
	if isAdmin {
		return fmt.Sprintf("%s/admins", *customerGroup.Path), nil
	}
	return fmt.Sprintf("%s/users", *customerGroup.Path), nil
}

func (c *UserContext) createUserAndSendEmail(clientId *string, email string, groupPath string) error {
	inviteUserError := &types.InviteUsersError{}
	user := gocloak.User{
		Email:    gocloak.StringP(email),
		Enabled:  gocloak.BoolP(true),
		Username: gocloak.StringP(email),
	}

	attributes := make(map[string][]string)
	attributes["approved"] = []string{"true"}
	user.Attributes = &attributes

	user.Groups = &[]string{groupPath}

	log.Info().Msg(fmt.Sprintf("Inviting user %s", *user.Email))
	search := gocloak.GetUsersParams{
		Username: &email,
	}
	users, err := c.client.GetUsers(c.ctx, c.token, types.KEYCLOAK_REALM, search)
	if err != nil {
		return errors.Wrapf(inviteUserError.InternalServerError(), "cannot find the user in realm %s",
			types.KEYCLOAK_REALM)
	}

	var userId string
	if len(users) == 0 {
		userId, err = c.client.CreateUser(c.ctx, c.token, types.KEYCLOAK_REALM, user)
		if err != nil {
			return errors.Wrapf(inviteUserError.InternalServerError(), "could not create user %s",
				user)
		}
		log.Debug().Msg(fmt.Sprintf("Created user with ID %s", userId))
	} else {
		userId = *users[0].ID
		log.Debug().Msg(fmt.Sprintf("Found existing user %s with ID %s", *user.Username, userId))
	}

	log.Debug().Msg(fmt.Sprintf("Sending invite email to %s", email))
	params := gocloak.ExecuteActionsEmail{}
	params.UserID = &userId
	params.ClientID = clientId
	// 30 minutes
	lifespan := types.EMAIL_LINK_DURATION_MINUTES * 60
	params.Lifespan = &lifespan
	params.Actions = &[]string{"UPDATE_PASSWORD"}

	err = c.client.ExecuteActionsEmail(c.ctx, c.token, c.realm, params)
	if err != nil {
		log.Err(err).Msg("Failed to invoke update password action email")
		return err
	}

	log.Debug().Msg(fmt.Sprintf("Sent invite email with duration of %d minutes", types.EMAIL_LINK_DURATION_MINUTES))
	return err
}

func getClientIdFromToken(tokenString string) (*string, error) {
	// The token is not verified or trusted, and it is supposed to be validated by keycloak security
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse the claims in the access token")
	}
	claim, err := extractClaimFromToken(claims, types.CLIENT_ID_CLAIM_NAME)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg(fmt.Sprintf("Found %s token claim as %s", types.ORG_ID_CLAIM_NAME, claim))
	return &claim, nil
}
