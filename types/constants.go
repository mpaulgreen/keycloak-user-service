package types

var (
	// Set a comma-delimited list of network origins (IPv4 addresses,
	// IPv4 CIDRs, IPv6 addresses or IPv6 CIDRs) from which to trust
	// request's headers that contain alternative client IP
	TRUSTED_PROXIES                    = "10.0.0.0/8"
	KEYCLOAK_BACKEND_URL               = "http://localhost:8080"
	CLIENT_ID                          = "admin-cli"
	KEYCLOAK_REALM                     = "master"
	KEYCLOAK_USERS_RESOURCE_URI        = "admin/realms/master/users"
	USER_SERVICE_TLS_CRT_PATH          = "./tls.crt"
	USER_SERVICE_TLS_KEY_PATH          = "./tls.key"
	KEYCLOAK_TLS_CRT_PATH              = "./tls.crt"
	KEYCLOAK_TLS_KEY_PATH              = "./tls.key"
	KEYCLOAK_CA_PATH                   = "./ca.crt"
	USER_SERVICE_PORT                  = ":8443"
	DISABLE_KEYCLOAK_CERT_VERIFICATION = "false"
	EMAIL_LINK_DURATION_MINUTES        = 30

	USERNAME_PARAM   = "username"
	ORG_ID_PARAM     = "org_id"
	USER_NAMES_PARAM = "usernames"
	EMAILS_PARAM     = "emails"
	USER_IDS_PARAM   = "user_ids"
	ID_PARAM         = "id"
	ACTIVATE_PARAM   = "activate"
	ISADMIN_PARAM    = "isAdmin"

	ORG_ID_CLAIM_NAME    = "organization.id"
	CLIENT_ID_CLAIM_NAME = "azp"

	AUTHORIZATION_HEADER = "Authorization"
	EFFECTIVE_TOKEN_KEY  = "token_key"

	APPROVED_ATTRIBUTE_NAME = "approved"
	ORG_ID_ATTRIBUTE        = "org_id"

	ORDER_BY_EMAIL    = "email"
	ORDER_BY_USERNAME = "username"
	ORDER_BY_MODIFIED = "modified"
	ORDER_BY_CREATED  = "created"

	ORDER_BY_DIR_ASC  = "asc"
	ORDER_BY_DIR_DESC = "desc"

	ORG_ADMIN_ATTRIBUTE = "org_admin"

	// http codes
	HTTP_CODE_BAD_REQUEST = "400"

	// Error messages
	ERR_NIL_HTTP_CLIENT_OR_REQUEST = "nil http request or http client object"

	RUN_ON_LOCAL           = "local"
	RUN_ON_OPENSHIFT_LOCAL = "openshift.local"
	RUN_ON_OPENSHIFT_DEV   = "openshift.dev"
	RUN_ON_DOCKER          = "docker"
	TRUE                   = true
	FALSE                  = false
	CORS_ALLOW_ORIGIN      = "https://stage.foo.redhat.com:1337,https://stage.foo.redhat.com"
)
