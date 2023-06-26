package user_routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	userhandlers "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/types"
	"net/http"
)

func GetUsersByUsersCriteria(c *gin.Context) {
	var findUsersCriteria types.FindUsersCriteria
	accessToken, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	if err = c.Bind(&findUsersCriteria); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	log.Debug().Msg(fmt.Sprintf("FindUsers Pagination called with %s: %+v\n", c.Request.RequestURI, findUsersCriteria))

	context, err := userhandlers.NewContext(accessToken)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to search users: %s", err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	usersList, code, err := context.FindUsers(findUsersCriteria)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to search users: %s", err))
		c.JSON(code, err)
		return
	}
	c.JSON(http.StatusOK, usersList)
}
