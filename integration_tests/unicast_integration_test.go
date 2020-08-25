package integration_tests

import (
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/uyuni/client"
)

func Test_Unicast(t *testing.T) {
	tt := []struct {
		name, call              string
		loginCredentials        struct{ username, password string }
		unicastResponseAnalizer func(unicastResponse interface{}) bool
		expectedError           string
	}{
		{
			name:                    "unicast.system.listSystems should succeed",
			call:                    "unicast.system.listSystems",
			loginCredentials:        struct{ username, password string }{"admin", "admin"},
			unicastResponseAnalizer: analizeListSystemsUnicastResponse,
		},
		{
			name:             "unkown method should fail",
			call:             "unicast.unkown.unkown",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			expectedError:    "request error: bad status code - 400",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			client := client.NewClient(10, 10)
			//login to Hub server
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAuthRelayMode", []interface{}{tc.loginCredentials.username, tc.loginCredentials.password})
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("Error occurred when executing login: %v", err)
			}
			hubSessionKey := loginResponse.(string)
			//login to peripheral server 1
			_, err = client.ExecuteCall(gatewayServerURL, "hub.attachToServers", []interface{}{hubSessionKey, []interface{}{peripheralServer1.id}})
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("Error occurred when executing attachToServers: %v", err)
			}
			//execute unicast call for peripheral server 1
			unicastResponse, err := client.ExecuteCall(gatewayServerURL, tc.call, []interface{}{hubSessionKey, peripheralServer1.id})
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("Error occurred when executing unicast call: %v", err)
			}
			if err == nil && (tc.expectedError != "" || !tc.unicastResponseAnalizer(unicastResponse)) {
				t.Fatalf("Expected and actual unicast responses don't match. Actual response is: %v", unicastResponse)
			}
			//logout
			_, err = client.ExecuteCall(gatewayServerURL, "hub.logout", []interface{}{hubSessionKey})
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("Error occurred when executing logout: %v", err)
			}
		})
	}
}

func analizeListSystemsUnicastResponse(unicastResponse interface{}) bool {
	minions := unicastResponse.([]interface{})
	if len(minions) != len(peripheralServer1.minions) {
		return false
	}
	for _, minion := range minions {
		minionMap := minion.(map[string]interface{})
		minionID := minionMap["id"].(int64)
		if !compareMinion(peripheralServer1.minions[minionID], minionMap) {
			return false
		}
	}
	return true
}
