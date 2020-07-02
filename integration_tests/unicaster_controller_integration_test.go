package integration_tests

import (
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

func Test_Unicast(t *testing.T) {
	tt := []struct {
		name, call             string
		loginCredentials       struct{ username, password string }
		analizeUnicastResponse func(unicastResponse interface{}) bool
		expectedError          string
	}{
		{
			name:                   "unicast.system.listSystems should succeed",
			call:                   "unicast.system.listSystems",
			loginCredentials:       struct{ username, password string }{"admin", "admin"},
			analizeUnicastResponse: analizeListSystemsUnicastResponse,
		},
		{
			name:                   "unkown method should fail",
			call:                   "unicast.unkown.unkown",
			loginCredentials:       struct{ username, password string }{"admin", "admin"},
			analizeUnicastResponse: analizeUnkonwMethodUnicastResponse,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAutoconnectMode", []interface{}{tc.loginCredentials.username, tc.loginCredentials.password})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			hubSessionKey := loginResponse.(map[string]interface{})["SessionKey"].(string)
			loggedInServerIDs := getLoggedInServerIDsFromLoginResponse(loginResponse)
			//execute unicast call
			unicastResponse, err := client.ExecuteCall(gatewayServerURL, tc.call, []interface{}{hubSessionKey, loggedInServerIDs[0]})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !tc.analizeUnicastResponse(unicastResponse) {
				t.Fatalf("Expected and actual values don't match. Actual value is: %v", unicastResponse)
			}
		})
	}
}

func analizeListSystemsUnicastResponse(unicastResponse interface{}) bool {
	// [map[id:1000010001 name:peripheral-server-1000010001-minion-2] map[id:1000010000 name:peripheral-server-1000010001-minion-1]]
	// systems := unicastResponse.([]interface{})

	// if len(systems) != 2 {
	// 	return false
	// }
	// for i, minionID := range systems {
	// 	if _, ok := peripheralServers[serverID.(int64)]; !ok {
	// 		return false
	// 	}
	// 	systems := successfulServerResponses[i].([]interface{})
	// 	for _, system := range systems {
	// 		systemMap := system.(map[string]interface{})
	// 		if !compareMinion(peripheralServers[serverID.(int64)].minions[systemMap["id"].(int64)], systemMap) {
	// 			return false
	// 		}
	// 	}
	// }
	return false
}

func analizeUnkonwMethodUnicastResponse(unicastResponse interface{}) bool {
	return false
}
