package integration_tests

import (
	"reflect"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

func Test_Unicast(t *testing.T) {
	tt := []struct {
		name             string
		loginCredentials struct{ username, password string }
		expectedResponse *controller.MulticastResponse
		expectedError    string
	}{
		{
			name:             "unicast.system.listSystems",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			expectedResponse: &controller.MulticastResponse{controller.MulticastStateResponse{}, controller.MulticastStateResponse{}},
		},
		// {methodName: "unicast.system.listUserSystems", args: []interface{}{"admin"}},
		// {methodName: "unicast.system.unknownmethod", args: []interface{}{"admin"}, output: "request error: bad status code - 400"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			peripheralServers := peripheralServers()
			initInfrastructure(peripheralServers, 8001, "admin", "admin")
			gatewayServerURL := "http://localhost:2830/hub/rpc/api"
			client := client.NewClient(10, 10)
			//login
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAutoconnectMode", []interface{}{"admin", "admin"})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			hubSessionKey := loginResponse.(map[string]interface{})["SessionKey"].(string)
			serverIDsSlice := loginResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})

			loggedInServerIDs := make([]int64, 0, len(serverIDsSlice))
			for _, serverID := range serverIDsSlice {
				loggedInServerIDs = append(loggedInServerIDs, serverID.(int64))
			}
			//execute multicast call
			systemsPerServer, err := client.ExecuteCall(gatewayServerURL, "unicast.system.listSystems", []interface{}{hubSessionKey, loggedInServerIDs[0]})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(systemsPerServer, tc.expectedResponse) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedResponse)
			}
		})
	}
}
