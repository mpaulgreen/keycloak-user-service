package user_handles

import (
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"keycloak-user-service/types"
	"net/http"
	strconv "strconv"

	"github.com/rs/zerolog/log"
)

func (c *CallContext) SetAdminStatus(userId string, isAdmin bool) (int, string) {
	log.Info().Msg("SetAdminStatus user with id: " + userId + " with value: " + strconv.FormatBool(isAdmin))

	user, err := c.client.GetUserByID(c.ctx, c.token, c.realm, userId)
	if err != nil {
		log.Err(err).Msg("Error retrieving user")
		return http.StatusInternalServerError, fmt.Sprintf("Error retrieving user with ID of [%s]!", userId)
	} else if user == nil {
		log.Error().Msg(fmt.Sprintf("Failed to retrieve user [%s], the caller token may lack access to fetch this user"))
		return http.StatusBadRequest, fmt.Sprintf("User with ID of [%s] not found!", userId)
	} else if isAdmin && !*user.Enabled {
		return http.StatusBadRequest, fmt.Sprintf("User with ID [%s] is not enabled and cannot be made an admin", userId)
	}

	orgId, err := c.getOrgIdFromToken(c.token)
	if err != nil {
		log.Err(err).Msg("Error getting org Id from token")
		return http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve %s from the access token: %s", types.ORG_ID_CLAIM_NAME, err.Error())
	}

	params := gocloak.GetGroupsParams{
		Q:                   gocloak.StringP(fmt.Sprintf("%s:%d", types.ORG_ID_ATTRIBUTE, *orgId)),
		Full:                gocloak.BoolP(true),
		BriefRepresentation: gocloak.BoolP(false),
	}

	groups, err := c.client.GetGroups(c.ctx, c.token, c.realm, params)
	if err != nil {
		log.Err(err).Msg("Error getting caller customer group")
		return http.StatusInternalServerError, fmt.Sprintf("Failed to load customer group with org Id [%s]: %s", *orgId, err.Error())
	} else if len(groups) == 0 {
		return http.StatusExpectationFailed, fmt.Sprintf("Cannot find customer group with %s attribute set to %d", types.ORG_ID_ATTRIBUTE, orgId)
	} else if len(groups) > 1 {
		return http.StatusExpectationFailed, fmt.Sprintf("Unexepctedly found %d groups with %s attribute set to %d", len(groups), types.ORG_ID_ATTRIBUTE, orgId)
	}

	customerGroup := groups[0]
	var adminsGroup, usersGroup *gocloak.Group
	for index, subGroup := range *customerGroup.SubGroups {
		if *subGroup.Name == "users" {
			usersGroup = &(*customerGroup.SubGroups)[index]
		} else if *subGroup.Name == "admins" {
			adminsGroup = &(*customerGroup.SubGroups)[index]
		}
	}
	if adminsGroup == nil {
		return http.StatusExpectationFailed, fmt.Sprintf("Failed to find an [admins] group under customer group %s", customerGroup.Name)
	} else if usersGroup == nil {
		return http.StatusExpectationFailed, fmt.Sprintf("Failed to find an [users] group under customer group %s", customerGroup.Name)
	}
	if isAdmin {
		err = c.client.AddUserToGroup(c.ctx, c.token, c.realm, userId, *adminsGroup.ID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to add user to admins group with Id [%s]: %s", adminsGroup.ID, err.Error())
		}
		err = c.client.DeleteUserFromGroup(c.ctx, c.token, c.realm, userId, *usersGroup.ID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to remove user from users group with Id [%s]: %s", usersGroup.ID, err.Error())
		}
		return http.StatusOK, fmt.Sprintf("Now user %s is an admin.", userId)
	} else {
		err = c.client.AddUserToGroup(c.ctx, c.token, c.realm, userId, *usersGroup.ID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to add user to users group with Id [%s]: %s", usersGroup.ID, err.Error())
		}
		err = c.client.DeleteUserFromGroup(c.ctx, c.token, c.realm, userId, *adminsGroup.ID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to remove user from admins group with Id [%s]: %s", adminsGroup.ID, err.Error())
		}
		return http.StatusOK, fmt.Sprintf("Now user %s is no longer an admin.", userId)
	}
}
