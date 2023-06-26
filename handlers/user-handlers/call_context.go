package user_handles

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"keycloak-user-service/client"
	"keycloak-user-service/types"
	"os"
	"strings"
)

type CallContext struct {
	client *gocloak.GoCloak
	token  string
	ctx    context.Context
	realm  string
}

// CAUTION: For use with unit tests only
var callCtxForUnitTests = CallContext{
	client: gocloak.NewClient(types.KEYCLOAK_BACKEND_URL),
	ctx:    context.Background(),
	realm:  types.KEYCLOAK_REALM,
}

// GetGoCloakClientForUnitTests For use with unit tests only
func GetGoCloakClientForUnitTests() *gocloak.GoCloak {
	return callCtxForUnitTests.client
}

func NewContext(accessToken string) (CallContext, error) {
	if run := os.Getenv("UNIT_TEST_RUN"); run != "" {
		log.Debug().Msg("Returning call context for unit tests")
		return callCtxForUnitTests, nil
	}

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
	}, nil
}
