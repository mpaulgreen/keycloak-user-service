package user_handles

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"keycloak-user-service/types"
	"strconv"
	"strings"
)

func (c *UserContext) getOrgIdFromToken(tokenString string) (*int, error) {
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

func (c *UserContext) populateAncestry(group types.GroupWrapper) error {
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
