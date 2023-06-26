package user_routes

import (
	"fmt"
	user_handles "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/types"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func InviteUser(c *gin.Context) {
	accessToken, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	var inviteUserData types.InviteUsers
	err = c.Bind(&inviteUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(inviteUserData)
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if inviteUserData.Emails == nil || len(inviteUserData.Emails) == 0 {
		log.Error().Msg("Must include emails of users to be invited")
		c.JSON(http.StatusBadRequest, errors.New("must include emails of users to be invited"))
		return
	}

	log.Debug().Msg(fmt.Sprintf("InviteUser payload: %v", inviteUserData))

	context, err := user_handles.NewContext(accessToken)
	if err != nil {
		err = &types.Error{Detail: "Failed to create a client", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	response, err := context.InviteUsers(inviteUserData)
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// At the moment just return success code created
	c.JSON(response.Code, response.Message)
}
