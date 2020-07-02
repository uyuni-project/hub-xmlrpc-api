package integration_tests

import (
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

func Test_Logout(t *testing.T) {
	tt := []struct {
		name                     string
		loginCredentials         struct{ username, password string }
		logoutParametersResolver func(sessionKey string) []interface{}
		expectedError            string
	}{
		{
			name:             "hub.logout should succeed",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			logoutParametersResolver: func(sessionKey string) []interface{} {
				return []interface{}{sessionKey}
			},
		},
		{
			name:             "hub.logout invalid SessionKey parameter should fail",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			logoutParametersResolver: func(sessionKey string) []interface{} {
				return []interface{}{"invalid_session_key"}
			},
			expectedError: "Authentication error: provided session key is invalid",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.login", []interface{}{tc.loginCredentials.username, tc.loginCredentials.password})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error executing login: %v", err)
			}
			hubSessionKey := loginResponse.(string)
			//logout
			_, err = client.ExecuteCall(gatewayServerURL, "hub.logout", tc.logoutParametersResolver(hubSessionKey))
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("Error executing logout: %v", err)
			}
		})
	}
}
