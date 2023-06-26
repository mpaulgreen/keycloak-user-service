package test

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	env "keycloak-user-service/env"
	"os"
)

func SetUpRouter() *gin.Engine {
	env.LoadEnvVars()
	log.Debug().Msg("Loaded environment variables")

	err := os.Setenv(UNIT_TEST_RUN, "true")
	if err != nil {
		log.Err(err).Msg("Error setting unit test environment variable")
	}
	router := gin.Default()
	return router
}
