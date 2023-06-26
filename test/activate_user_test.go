package test

import (
	user_handles "keycloak-user-service/handlers/user-handlers"
	"keycloak-user-service/routes"
	userroutes "keycloak-user-service/routes/user-routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestActivateUserNoId(t *testing.T) {
	router := routes.Router{Server: SetUpRouter()}

	router.HandleRoute(router.Server.PUT, "/user/:id/activate/:activate", userroutes.ActivateUser)

	req := GetNewHttpRequest(http.MethodPut, "/user/"+EMPTY_PARAM_WITH_SPACE+"/activate/true", nil)
	w := httptest.NewRecorder()
	router.Server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestActivateUserNoActivateValue(t *testing.T) {
	router := routes.Router{Server: SetUpRouter()}
	router.HandleRoute(router.Server.PUT, "/user/:id/activate/:activate", userroutes.ActivateUser)
	req := GetNewHttpRequest(http.MethodPut, "/user/"+KEYCLOAK_ACTIVATE_USER_ID+"/activate/"+EMPTY_PARAM_WITH_SPACE, nil)
	w := httptest.NewRecorder()
	router.Server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

//func TestActivateUserTrue(t *testing.T) { TODO fix test, gocloak returns an empty user when trying to load by ID
//	httpmock.ActivateNonDefault(user_handles.GetGoCloakClientForUnitTests().RestyClient().GetClient())
//	defer httpmock.DeactivateAndReset()
//	setupHttpMockForActivateUser()
//	router := routes.Router{Server: SetUpRouter()}
//	router.HandleRoute(router.Server.PUT, "/user/:id/activate/:activate", userroutes.ActivateUser)
//	req := GetNewHttpRequest(http.MethodPut, "/user/"+KEYCLOAK_ACTIVATE_USER_ID+"/activate/true", nil)
//	w := httptest.NewRecorder()
//	router.Server.ServeHTTP(w, req)
//	assert.Equal(t, http.StatusNoContent, w.Code)
//}

func TestActivateUserFalse(t *testing.T) {
	httpmock.ActivateNonDefault(user_handles.GetGoCloakClientForUnitTests().RestyClient().GetClient())
	defer httpmock.DeactivateAndReset()

	setupHttpMockForActivateUser()
	router := routes.Router{Server: SetUpRouter()}

	router.HandleRoute(router.Server.PUT, "/user/:id/activate/:activate", userroutes.ActivateUser)
	req := GetNewHttpRequest(http.MethodPut, "/user/"+KEYCLOAK_ACTIVATE_USER_ID+"/activate/false", nil)
	w := httptest.NewRecorder()
	router.Server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func setupHttpMockForActivateUser() {
	httpmock.RegisterResponder("PUT", KEYCLOAK_ACTIVATE_USER_URL,
		httpmock.NewStringResponder(204, "[]"))
	httpmock.RegisterResponder("PUT", KEYCLOAK_USERS_PATH,
		httpmock.NewStringResponder(204, "[]"))
	httpmock.RegisterResponder("GET", KEYCLOAK_ACTIVATE_USER_GET_GROUPS_URL,
		httpmock.NewStringResponder(200, KEYCLOAK_USER_DATA1))
	httpmock.RegisterResponder("GET", KEYCLOAK_ACTIVATE_USER_GET_USER_URL,
		httpmock.NewStringResponder(200, KEYCLOAK_USER_DATA1))
}
