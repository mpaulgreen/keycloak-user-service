package main

import (
	"github.com/rs/zerolog/log"
	"net/http"
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
	engine.GET("/health/ready", isHealthy)
	engine.GET("/health/live", isHealthy)
	router.HandleRoute(router.Server.GET, "/users", userroutes.GetUsersByUsersCriteria)
	router.HandleRoute(router.Server.PUT, "/user/:id/activate/:activate", userroutes.ActivateUser)
	router.HandleRoute(router.Server.PUT, "/change-users-status", userroutes.ChangeUsersStatus)
	router.HandleRoute(router.Server.PUT, "/user/invite", userroutes.InviteUser)
	router.HandleRoute(router.Server.PUT, "/user/:id/admin/:isAdmin", userroutes.SetAdminStatus)

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

func isHealthy(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
