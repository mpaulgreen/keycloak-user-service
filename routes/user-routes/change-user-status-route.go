package user_routes

import (
	"github.com/gin-gonic/gin"
	userhandlers "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/types"
	"net/http"
)

func ChangeUsersStatus(c *gin.Context) {
	var reqBody types.ChangeUsersStatusDTO
	accessToken, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	if err = c.BindJSON(&reqBody); err != nil {
		errorResponse := &types.Error{Detail: "Missing request body or request body is not formed correctly as expected.", Status: types.HTTP_CODE_BAD_REQUEST}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var response []types.ChangeUsersStatus

	for _, user := range reqBody.ChangeUsersStatus {
		var respUser types.ChangeUsersStatus
		respUser.UserId = user.UserId
		respUser.IsActive = user.IsActive

		context, err := userhandlers.NewContext(accessToken)
		if err != nil {
			err = &types.Error{Detail: "Failed to create a client", Status: types.HTTP_CODE_BAD_REQUEST}
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if err = context.ActivateUser(user.UserId, user.IsActive); err != nil {
			respUser.Status = http.StatusBadRequest
			respUser.Msg = "Error changing status of the user with id=" + user.UserId
			return
		} else {
			respUser.Status = http.StatusNoContent
			respUser.Msg = "Successful"
		}
		response = append(response, respUser)
	}

	c.JSON(http.StatusAccepted, &response)

	return

}
