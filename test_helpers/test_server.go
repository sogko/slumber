package test_helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"net/http"
	"net/http/httptest"
)

type TestServerOptions struct {
	RequestAcceptHeader string
	ServerConfig        *server.Config
	PrivateSigningKey   []byte
	PublicSigningKey    []byte
}

type TestServer struct {
	Options        *TestServerOptions
	TokenAuthority domain.ITokenAuthority
}

type AuthOptions struct {
	APIUser *domain.User
	Token   string
}

func NewTestServer(options *TestServerOptions) *TestServer {

	ta := middlewares.NewTokenAuthority(&middlewares.TokenAuthorityOptions{
		PrivateSigningKey: options.PrivateSigningKey,
		PublicSigningKey:  options.PublicSigningKey,
	})

	ts := TestServer{options, ta}
	return &ts

}
func (ts *TestServer) Request(recorder *httptest.ResponseRecorder, method string, urlStr string, body interface{}, targetResponse interface{}, authOptions *AuthOptions) {

	var s *server.Server
	var request *http.Request

	// request for version 0.0
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		request, _ = http.NewRequest(method, urlStr, bytes.NewReader(jsonBytes))
	} else {
		request, _ = http.NewRequest(method, urlStr, nil)
	}
	// set API version through accept header
	request.Header.Set("Accept", ts.Options.RequestAcceptHeader)

	if authOptions == nil {
		authOptions = &AuthOptions{nil, ""}
	}
	if authOptions.APIUser != nil {
		// set Authorization header
		var rolesString []string
		for _, role := range authOptions.APIUser.Roles {
			rolesString = append(rolesString, string(role))
		}
		token, _ := ts.TokenAuthority.CreateNewSessionToken(&domain.TokenClaims{
			UserID:   authOptions.APIUser.ID.Hex(),
			Username: authOptions.APIUser.Username,
			Status:   authOptions.APIUser.Status,
			Roles:    rolesString,
		})
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	} else {
		if authOptions.Token != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", authOptions.Token))
		}
	}

	// init server
	s = server.NewServer(ts.Options.ServerConfig).SetupRoutes()

	// serve request
	s.ServeHTTP(recorder, request)
	DecodeResponseToType(recorder, &targetResponse)

}
