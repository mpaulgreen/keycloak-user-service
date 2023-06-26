package user_routes

import (
	"errors"
	"keycloak-user-service/types"

	"github.com/gin-gonic/gin"
)

func GetAccessToken(c *gin.Context) (string, error) {
	accessToken := c.MustGet(types.EFFECTIVE_TOKEN_KEY).(string)
	if len(accessToken) < 10 {
		//Bearer prefix itself is 7 characters
		return "", errors.New("no access token of proper length was found in the request")
	}
	return accessToken, nil
}
