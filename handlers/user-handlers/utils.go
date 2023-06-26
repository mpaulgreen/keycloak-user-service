package user_handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"keycloak-user-service/types"
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

func extractClaimFromToken(claims jwt.MapClaims, claimName string) (string, error) {
	claim, ok := claims[claimName]
	if !ok {
		subkeys := strings.Split(claimName, ".")
		if len(subkeys) > 1 {
			subkey := subkeys[1]
			key := subkeys[0]
			if claim, ok = claims[key]; ok {
				if dict, ok := claim.(map[string]interface{}); ok {
					if value, ok := dict[subkey]; ok {
						if valueStr, ok := value.(string); ok {
							return valueStr, nil
						}
						return "", fmt.Errorf("value of parameter %s in the access token from the %s claim is not a string", subkey, key)
					}
					return "", fmt.Errorf("could not find the %s parameter in the access token from the %s claim", subkey, key)
				}
				log.Debug().Msg(fmt.Sprintf("Claim %s in the access token is not of the expected type, but %s", key, reflect.TypeOf(claim)))
				return "", fmt.Errorf("claim %s in the access token is not of the expected type", key)
			}
			log.Debug().Msg(fmt.Sprintf("Missing claim %s in the access token, found: %s", key, maps.Keys(claims)))
			return "", fmt.Errorf("missing claim %s in the access token", key)
		}
		return "", fmt.Errorf("could not find the %s parameter in the access token", claimName)
	}
	if claimStr, ok := claim.(string); ok {
		return claimStr, nil
	}
	return "", fmt.Errorf("value of the %s claim in the access token is not a string", claimName)
}

func GetAccessToken(c *gin.Context) (string, error) {
	accessToken := c.MustGet(types.EFFECTIVE_TOKEN_KEY).(string)
	if len(accessToken) < 10 {
		//Bearer prefix itself is 7 characters
		return "", errors.New("no access token of proper length was found in the request")
	}
	return accessToken, nil
}
