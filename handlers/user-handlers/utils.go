package user_handles

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

func handleNilString(str *string) string {
	if str == nil {
		return ""
	} else {
		return *str
	}
}

func handleNilBool(boolean *bool) bool {
	if boolean == nil {
		return false
	} else {
		return *boolean
	}
}

func containsAttributeValue(attributes *map[string][]string, key string, expectedValue string) bool {
	if attributes != nil {
		for _, value := range (*attributes)[key] {
			if value == expectedValue {
				return true
			} else {
				log.Debug().Msg(fmt.Sprintf("Attributes map contains attribute %s but with different value %s from the expected %s",
					key, value, expectedValue))
			}
		}
	}
	return false
}

func getSingleAttributeValue(attributes *map[string][]string, key string) (*string, error) {
	values, found := (*attributes)[key]
	if !found {
		return nil, nil
	} else if len(values) > 1 {
		err := errors.New(fmt.Sprintf("More than one value found for the %s attribute, aborting!", key))
		log.Err(err).Msg("Found more than the single expected attribute")
		return nil, err
	} else {
		return &values[0], nil
	}
}

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

func hasContent(array []string) bool {
	return array != nil && len(array) > 0
}

func equalsStringArrays(a, b []string) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else if len(a) != len(b) {
		return false
	} else {
		for index := 0; index < len(a); index++ {
			if a[index] != b[index] {
				return false
			}
		}
	}
	return true
}
