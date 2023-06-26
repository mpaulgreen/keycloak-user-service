package routes

import (
	"keycloak-user-service/types"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Router struct {
	Server *gin.Engine
}

func (router *Router) HandleRoute(httpMethod func(string, ...gin.HandlerFunc) gin.IRoutes, path string, handlers ...gin.HandlerFunc) {
	handlerChain := []gin.HandlerFunc{addCORSHeaders, getCallerToken}
	handlerChain = append(handlerChain, handlers...)
	httpMethod(path, handlerChain...)
	router.Server.OPTIONS(path, addCORSHeaders)
}

func addCORSHeaders(c *gin.Context) {
	log.Info().Msg("Adding CORS headers")
	c.Header("Access-Control-Allow-Origin", types.CORS_ALLOW_ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS,HEAD")
	c.Header("Access-Control-Allow-Headers", "Origin,Content-type,Accept,Access-Control-Allow-Origin,Authorization")
	c.Header("Access-Control-Max-Age", "3600")
	c.Next()
}

func getCallerToken(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) < 10 {
		//Bearer prefix itself is 7 characters
		log.Error().Msg("Caller did not provide a token: [" + accessToken + "]")
	} else {
		tokenParts := strings.Split(accessToken, ".")
		if len(tokenParts) == 3 {
			log.Debug().Msg("JWT payload: [" + tokenParts[1] + "]")
		} else {
			log.Debug().Msg("Non-standard token: [" + accessToken + "]")
		}
	}
	c.Set(types.EFFECTIVE_TOKEN_KEY, accessToken)
	c.Next()
}
