package user_handles

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"keycloak-user-service/client"
	"keycloak-user-service/types"
	"strings"
)

type UserContext struct {
	client *gocloak.GoCloak
	token  string
	ctx    context.Context
	realm  string
	ginCtx *gin.Context
}

func (uc *UserContext) NewContext(accessToken string, ginCtx *gin.Context) error {
	kclient, err := client.NewClient()
	if err != nil {
		return err
	}
	uc.client = kclient
	uc.token = strings.Replace(accessToken, "Bearer ", "", 1)
	uc.ctx = context.Background()
	uc.realm = types.KEYCLOAK_REALM
	uc.ginCtx = ginCtx
	return nil
}
