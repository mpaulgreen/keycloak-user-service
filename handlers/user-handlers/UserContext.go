package user_handles

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"keycloak-user-service/client"
	"keycloak-user-service/types"
	"strings"
)

type CallContext struct {
	client *gocloak.GoCloak
	token  string
	ctx    context.Context
	realm  string
	ginCtx *gin.Context
}

func NewContext(accessToken string, c *gin.Context) (CallContext, error) {
	kclient, err := client.NewClient()
	if err != nil {
		return CallContext{}, err
	}
	ctx := context.Background()
	return CallContext{
		client: kclient,
		token:  strings.Replace(accessToken, "Bearer ", "", 1),
		ctx:    ctx,
		realm:  types.KEYCLOAK_REALM,
		ginCtx: c,
	}, nil
}
