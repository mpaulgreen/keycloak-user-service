package user_routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	userhandlers "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/types"
	"net/http"
	"strconv"
)

func SetAdminStatus(c *gin.Context) {
	accessToken, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	isAdmin, err := strconv.ParseBool(c.Param(types.ISADMIN_PARAM))
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Provided isAdmin property of %s could not be parsed to true or false", c.Param(types.ISADMIN_PARAM)))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	context, err := userhandlers.NewContext(accessToken)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to set admin status: %s", err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	code, msg := context.SetAdminStatus(c.Param(types.ID_PARAM), isAdmin)
	c.JSON(code, gin.H{
		"code":   code,
		"detail": msg,
		"status": http.StatusText(code),
	})
}
