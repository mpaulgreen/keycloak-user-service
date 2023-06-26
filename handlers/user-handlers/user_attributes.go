package user_handles

import (
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"keycloak-user-service/types"
	"strconv"
	"strings"
)

// For a given userId, load all user groups, and calculate their inherited attributes
// Reconcile attributes from separate groups and return an error in case of conflicting values
func (c *CallContext) effectiveAttributes(user *gocloak.User) (*map[string][]string, error) {
	if user == nil {
		return nil, errors.New("User object is nil")
	}
	groups, err := c.groupMembership(*user.ID)
	if err != nil {
		return nil, err
	}
	attrs, err := reconcileAttributes(groups)
	if err != nil {
		return nil, err
	}
	if user.Attributes != nil {
		for key, value := range *user.Attributes {
			(*attrs)[key] = value
		}
	}
	return attrs, nil
}

func (c *CallContext) getOrgIdFromToken(tokenString string) (*int, error) {
	// The token is not verified or trusted, and it is supposed to be validated by keycloak security
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse the claims in the access token")
	}
	claim, err := extractClaimFromToken(claims, types.ORG_ID_CLAIM_NAME)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg(fmt.Sprintf("Found %s token claim as %s", types.ORG_ID_CLAIM_NAME, claim))
	oId, err := strconv.Atoi(claim)
	if err != nil {
		return nil, fmt.Errorf("the value of organization id parameter in the access token is invalid")
	}
	orgId := new(int)
	*orgId = oId
	return orgId, nil
}

func (c *CallContext) groupMembership(userId string) ([]types.GroupWrapper, error) {
	groups, err := c.client.GetUserGroups(c.ctx, c.token, c.realm, userId, gocloak.GetGroupsParams{
		BriefRepresentation: &types.FALSE,
	})
	if err != nil {
		return nil, err
	}

	var directGroups []types.GroupWrapper
	for _, group := range groups {
		loadedGroup := types.WrapGroup(group)
		err = c.populateAncestry(loadedGroup)
		if err != nil {
			return nil, err
		}
		directGroups = append(directGroups, loadedGroup)
	}
	return directGroups, nil
}

func (c *CallContext) getCurrentUserGroups() ([]types.GroupWrapper, error) {
	userInfo, err := c.client.GetUserInfo(c.ctx, c.token, c.realm)
	if err != nil {
		return nil, err
	}

	return c.groupMembership(*userInfo.Sub)
}

func (c *CallContext) populateAncestry(group types.GroupWrapper) error {
	path := *group.Group().Path
	name := fmt.Sprintf("/%s", *group.Group().Name)
	if len(path) > len(name) {
		//Group has parent(s), but the only way to find them is using the path.
		//While slash is the separator and is also allowed in group name, this is not an issue since the full path loads the group,
		// and then the name is removed to go up the hierarchy, along with the false separator
		remainingPath := strings.TrimSuffix(path, name)
		parent, err := c.client.GetGroupByPath(c.ctx, c.token, c.realm, remainingPath)
		if err != nil {
			return err
		}
		parentWrapper := group.SetParent(parent)
		err = c.populateAncestry(parentWrapper)
		if err != nil {
			return err
		}
	}
	return nil
}

func reconcileAttributes(groups []types.GroupWrapper) (*map[string][]string, error) {
	attrs := make(map[string][]string)
	for _, group := range groups {
		for key, values := range *group.InheritedAttributes() {
			existing, found := attrs[key]
			if found && !equalsStringArrays(values, existing) {
				msg := fmt.Sprintf("Conflicting values found for attribute %s in group %s from what was seen in another group", key, *group.Group().Name)
				return nil, errors.New(msg)
			} else {
				attrs[key] = values
			}
		}
	}
	return &attrs, nil
}

func orgIdAttribute(attrs *map[string][]string) (*string, error) {
	if attrs != nil {
		orgId, err := getSingleAttributeValue(attrs, types.ORG_ID_ATTRIBUTE)
		if err != nil {
			return nil, err
		}
		return orgId, nil
	}
	return nil, nil
}

func orgAdminAttribute(attrs *map[string][]string) (bool, error) {
	if attrs != nil {
		orgAdmin, err := getSingleAttributeValue(attrs, types.ORG_ADMIN_ATTRIBUTE)
		if err != nil {
			return false, err
		}
		if handleNilString(orgAdmin) == "true" {
			return true, nil
		}
	}
	return false, nil
}

func approvedAttribute(attrs *map[string][]string) (bool, error) {
	if attrs != nil {
		approved, err := getSingleAttributeValue(attrs, types.APPROVED_ATTRIBUTE_NAME)
		if err != nil {
			return false, err
		}
		if handleNilString(approved) == "true" {
			return true, nil
		}
	}
	return false, nil
}
