package user_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	userHandlers "keycloak-user-service/handlers/user-handlers"
	"net/http"
)

func InviteUser(c *gin.Context) {
	accessToken, err := GetAccessToken(c)
	if err != nil {
		log.Error().Msg("system error: Error in parsing token")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	uc := userHandlers.UserContext{}
	err = uc.NewContext(accessToken, c)
	if err != nil {
		log.Error().Msg("system error: failed to create client")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	uc.InviteUsers()
}
