package env

import (
	"fmt"
	"keycloak-user-service/types"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// LoadEnvVars Loads environment variables
func LoadEnvVars() {
	types.CLIENT_ID = os.Getenv("CLIENT_ID")
	types.USER_SERVICE_PORT = os.Getenv("USER_SERVICE_PORT")

	if len(types.USER_SERVICE_PORT) > 0 && !strings.HasPrefix(types.USER_SERVICE_PORT, ":") {
		types.USER_SERVICE_PORT = ":" + types.USER_SERVICE_PORT
	}

	types.KEYCLOAK_BACKEND_URL = os.Getenv("KEYCLOAK_BACKEND_URL")
	types.KEYCLOAK_REALM = os.Getenv("KEYCLOAK_REALM")
	types.KEYCLOAK_USERS_RESOURCE_URI = os.Getenv("KEYCLOAK_USERS_RESOURCE_URI")
	types.KEYCLOAK_USERS_RESOURCE_URI = os.ExpandEnv(types.KEYCLOAK_USERS_RESOURCE_URI)
	types.DISABLE_KEYCLOAK_CERT_VERIFICATION = os.Getenv("DISABLE_KEYCLOAK_CERT_VERIFICATION")
	types.USER_SERVICE_TLS_CRT_PATH = os.Getenv("USER_SERVICE_TLS_CRT_PATH")
	types.USER_SERVICE_TLS_KEY_PATH = os.Getenv("USER_SERVICE_TLS_KEY_PATH")
	types.KEYCLOAK_TLS_CRT_PATH = os.Getenv("KEYCLOAK_TLS_CRT_PATH")
	types.KEYCLOAK_TLS_KEY_PATH = os.Getenv("KEYCLOAK_TLS_KEY_PATH")
	types.KEYCLOAK_CA_PATH = os.Getenv("KEYCLOAK_CA_PATH")
	types.CORS_ALLOW_ORIGIN = os.Getenv("CORS_ALLOW_ORIGIN")
	if value, ok := os.LookupEnv("TRUSTED_PROXIES"); ok {
		types.TRUSTED_PROXIES = value
	}
	if orgIdClaimName, ok := os.LookupEnv("ORG_ID_CLAIM_NAME"); ok {
		types.ORG_ID_CLAIM_NAME = orgIdClaimName
	}
	if emailLinkDurationMinutesValue, ok := os.LookupEnv("EMAIL_LINK_DURATION_MINUTES"); ok {
		emailLinkDurationMinutes, err := strconv.Atoi(emailLinkDurationMinutesValue)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Skipping non-numeric value passed to EMAIL_LINK_DURATION_MINUTES variable: %s", emailLinkDurationMinutesValue))
		} else {
			types.EMAIL_LINK_DURATION_MINUTES = emailLinkDurationMinutes
		}
	}
}
