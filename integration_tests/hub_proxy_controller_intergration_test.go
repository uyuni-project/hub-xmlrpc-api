package integration_tests

import (
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

func Test_ProxyCallToHub(t *testing.T) {
	tt := []struct {
		name, call                     string
		loginCredentials               struct{ username, password string }
		proxyCallToHubResponseAnalizer func(multicastResponse interface{}) bool
		expectedError                  string
	}{
		{
			name:                           "system.listSystems should succeed",
			call:                           "system.listSystems",
			loginCredentials:               struct{ username, password string }{"admin", "admin"},
			proxyCallToHubResponseAnalizer: analizeListSystemsProxyCallToHubResponse,
		},
		{
			name:             "unkown method should fail",
			call:             "unkown.unkown",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			expectedError:    "request error: bad status code - 400",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.login", []interface{}{tc.loginCredentials.username, tc.loginCredentials.password})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error occurred when executing login: %v", err)
			}
			hubSessionKey := loginResponse.(string)
			//execute call
			callResponse, err := client.ExecuteCall(gatewayServerURL, tc.call, []interface{}{hubSessionKey})
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("Error occurred when executing call: %v", err)
			}
			if err == nil && !tc.proxyCallToHubResponseAnalizer(callResponse) {
				t.Fatalf("Expected and actual responses don't match. Actual response is: %v", callResponse)
			}
			//logout
			_, err = client.ExecuteCall(gatewayServerURL, "hub.logout", []interface{}{hubSessionKey})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error occurred when executing logout: %v", err)
			}
		})
	}
}

func analizeListSystemsProxyCallToHubResponse(callResponse interface{}) bool {
	minions := callResponse.([]interface{})
	if len(minions) != len(peripheralServers) {
		return false
	}
	for _, minion := range minions {
		minionMap := minion.(map[string]interface{})
		minionID := minionMap["id"].(int64)
		if !compareMinion(peripheralServers[minionID], minionMap) {
			return false
		}
	}
	return true
}
