package integration_tests

import (
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/uyuni/client"
)

func Test_Multicast(t *testing.T) {
	tt := []struct {
		name, call                string
		loginCredentials          struct{ username, password string }
		multicastResponseAnalizer func(multicastResponse interface{}) bool
		expectedError             string
	}{
		{
			name:                      "multicast.system.listSystems should succeed",
			call:                      "multicast.system.listSystems",
			loginCredentials:          struct{ username, password string }{"admin", "admin"},
			multicastResponseAnalizer: analizeListSystemsMulticastResponse,
		},
		{
			name:                      "unkown method should fail",
			call:                      "multicast.unkown.unkown",
			loginCredentials:          struct{ username, password string }{"admin", "admin"},
			multicastResponseAnalizer: analizeUnkonwMethodMulticastResponse,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAutoconnectMode", []interface{}{tc.loginCredentials.username, tc.loginCredentials.password})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error occurred when executing login: %v", err)
			}
			hubSessionKey := loginResponse.(map[string]interface{})["SessionKey"].(string)
			loggedInServerIDs := getLoggedInServerIDsFromLoginResponse(loginResponse)
			//execute multicast call
			multicastResponse, err := client.ExecuteCall(gatewayServerURL, tc.call, []interface{}{hubSessionKey, loggedInServerIDs})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error occurred when executing multicast call: %v", err)
			}
			if err == nil && !tc.multicastResponseAnalizer(multicastResponse) {
				t.Fatalf("Expected and actual multicast responses don't match. Actual response is: %v", multicastResponse)
			}
			//logout
			_, err = client.ExecuteCall(gatewayServerURL, "hub.logout", []interface{}{hubSessionKey})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error occurred when executing logout: %v", err)
			}
		})
	}
}

func getLoggedInServerIDsFromLoginResponse(loginResponse interface{}) []int64 {
	serverIDsSlice := loginResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})
	loggedInServerIDs := make([]int64, 0, len(serverIDsSlice))
	for _, serverID := range serverIDsSlice {
		loggedInServerIDs = append(loggedInServerIDs, serverID.(int64))
	}
	return loggedInServerIDs
}

func analizeListSystemsMulticastResponse(multicastResponse interface{}) bool {
	failedServersResponses := multicastResponse.(map[string]interface{})["Failed"].(map[string]interface{})["Responses"].([]interface{})
	failedServerIDs := multicastResponse.(map[string]interface{})["Failed"].(map[string]interface{})["ServerIds"].([]interface{})

	if len(failedServersResponses) != 0 || len(failedServerIDs) != 0 {
		return false
	}

	successfulServerResponses := multicastResponse.(map[string]interface{})["Successful"].(map[string]interface{})["Responses"].([]interface{})
	successfulServerIDs := multicastResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})

	if len(successfulServerIDs) != len(peripheralServers) || len(successfulServerResponses) != len(peripheralServers) {
		return false
	}
	for i, serverID := range successfulServerIDs {
		if _, ok := peripheralServers[serverID.(int64)]; !ok {
			return false
		}
		systems := successfulServerResponses[i].([]interface{})
		for _, system := range systems {
			systemMap := system.(map[string]interface{})
			if !compareMinion(peripheralServers[serverID.(int64)].minions[systemMap["id"].(int64)], systemMap) {
				return false
			}
		}
	}
	return true
}

func compareMinion(systemInfo SystemInfo, systemMap map[string]interface{}) bool {
	if systemID, ok := systemMap["id"]; !(ok && systemID.(int64) == systemInfo.id) {
		return false
	}
	if systemName, ok := systemMap["name"]; !(ok && systemName.(string) == systemInfo.name) {
		return false
	}
	return true
}

func analizeUnkonwMethodMulticastResponse(multicastResponse interface{}) bool {
	successfulServerResponses := multicastResponse.(map[string]interface{})["Successful"].(map[string]interface{})["Responses"].([]interface{})
	successfulServerIDs := multicastResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})

	if len(successfulServerResponses) != 0 || len(successfulServerIDs) != 0 {
		return false
	}

	failedServerResponses := multicastResponse.(map[string]interface{})["Failed"].(map[string]interface{})["Responses"].([]interface{})
	failedServerIDs := multicastResponse.(map[string]interface{})["Failed"].(map[string]interface{})["ServerIds"].([]interface{})

	if len(failedServerIDs) != len(peripheralServers) || len(failedServerResponses) != len(peripheralServers) {
		return false
	}
	for _, serverID := range failedServerIDs {
		if _, ok := peripheralServers[serverID.(int64)]; !ok {
			return false
		}
		for _, failedServerResponse := range failedServerResponses {
			if "request error: bad status code - 400" != failedServerResponse.(string) {
				return false
			}
		}
	}
	return true
}
