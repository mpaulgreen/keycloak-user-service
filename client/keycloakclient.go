package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"keycloak-user-service/types"
	"log"
	"net/http"
	"strconv"

	"github.com/Nerzal/gocloak/v13"
)

func NewClient() (*gocloak.GoCloak, error) {
	client := gocloak.NewClient(types.KEYCLOAK_BACKEND_URL)
	// restyClient.FormData.Add(types.AUTHORIZATION_HEADER, accessToken)

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = &x509.CertPool{}
	}

	restyClient := client.RestyClient()
	disableTlsCertVerification, _ := strconv.ParseBool(types.DISABLE_KEYCLOAK_CERT_VERIFICATION)
	if disableTlsCertVerification {
		restyClient.SetTLSClientConfig(&tls.Config{
			RootCAs:            rootCAs,
			InsecureSkipVerify: true,
		})
		return client, nil
	}

	// Read in the cert file
	caCerts, err := ioutil.ReadFile(types.KEYCLOAK_CA_PATH)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", types.KEYCLOAK_CA_PATH, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(caCerts); !ok {
		log.Println("No certs appended, using system certs only")
	}

	cert, err := tls.LoadX509KeyPair(types.KEYCLOAK_TLS_CRT_PATH, types.KEYCLOAK_TLS_KEY_PATH)
	if err != nil {
		return nil, err
	}

	restyClient.SetTLSClientConfig(&tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	})

	return client, nil
}

func ToError(err error, errorMessage string) types.Error {
	if apiError, ok := err.(*gocloak.APIError); ok {
		return types.Error{
			Code:   apiError.Code,
			Status: http.StatusText(apiError.Code),
			Detail: apiError.Message,
		}
	}
	if error, ok := err.(types.Error); ok {
		return error
	}

	return types.InternalError(err, errorMessage)
}
