package main

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"

	"keycloak-user-service/env"
	"keycloak-user-service/routes"
	userroutes "keycloak-user-service/routes/user-routes"
	"keycloak-user-service/types"

	"github.com/gin-gonic/gin"
)

func main() {
	env.LoadEnvVars()

	engine := gin.New()
	err := engine.SetTrustedProxies(strings.Split(types.TRUSTED_PROXIES, ","))
	if err != nil {
		log.Err(err).Msg("Failed to set trusted proxies.")
		panic(err)
	}
	router := &routes.Router{Server: engine}
	router.HandleRoute(router.Server.PUT, "/user/invite", userroutes.InviteUser)

	disableTlsCertVerification, _ := strconv.ParseBool(types.DISABLE_KEYCLOAK_CERT_VERIFICATION)
	if disableTlsCertVerification {
		err := router.Server.Run(types.USER_SERVICE_PORT)
		if err != nil {
			panic(err)
		}
	}

	err = router.Server.RunTLS(types.USER_SERVICE_PORT, types.USER_SERVICE_TLS_CRT_PATH, types.USER_SERVICE_TLS_KEY_PATH)
	if err != nil {
		panic(err)
	}
}
