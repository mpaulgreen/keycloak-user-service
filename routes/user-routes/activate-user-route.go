package user_routes

import (
	userhandlers "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ActivateUser(c *gin.Context) {
	id := strings.TrimSpace(c.Param(types.ID_PARAM))
	activate := strings.TrimSpace(c.Param(types.ACTIVATE_PARAM))
	accessToken, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	if id == "" {
		err = &types.Error{Detail: "Missing id parameter to activate user", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if activate == "" {
		err = &types.Error{Detail: "Missing activate parameter to activate user", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err != nil {
		err = &types.Error{Detail: "Invalid/missing user access token", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	activateBool := false
	if activate == "true" || activate == "1" {
		activateBool = true
	}

	context, err := userhandlers.NewContext(accessToken)
	if err != nil {
		err = &types.Error{Detail: "Failed to create a client", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err = context.ActivateUser(id, activateBool); err != nil {
		log.Err(err).Msg("Error activating user")
		c.JSON(http.StatusBadRequest, err)
		return
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
