package integration_tests

import (
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

// func init() {
// 	main.main()
// 	cmd := exec.Command("go", "run", "/home/marcelo/go/src/github.com/uyuni-project/hub-xmlrpc-api/hub_api_gateway.go")
// 	_, err := cmd.Output()
// 	if err != nil {
// 		log.Fatalf("cmd.Run() failed with %s\n", err)
// 	}
// }
func Test_Multicast(t *testing.T) {
	tt := []struct {
		name, username, password, call string
		expectedError                  string
	}{
		{
			name:     "multicast.system.listSystems",
			call:     "multicast.system.listSystems",
			username: "admin",
			password: "admin",
		},
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
			multicastResponse, err := client.ExecuteCall(gatewayServerURL, "multicast.system.listSystems", []interface{}{hubSessionKey, loggedInServerIDs})
			if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !analizeListSystemsResponse(peripheralServers, multicastResponse) {
				t.Fatalf("Expected and actual values don't match. Actual value is: %v, expected value is: %v", multicastResponse, peripheralServers)
			}
		})
	}
}

func analizeListSystemsResponse(peripheralServers map[int64]SystemInfo, listSystemsResponse interface{}) bool {
	failedServersResponses := listSystemsResponse.(map[string]interface{})["Failed"].(map[string]interface{})["Responses"].([]interface{})
	failedServerIDs := listSystemsResponse.(map[string]interface{})["Failed"].(map[string]interface{})["ServerIds"].([]interface{})

	if len(failedServersResponses) != 0 || len(failedServerIDs) != 0 {
		return false
	}

	successfulServerResponses := listSystemsResponse.(map[string]interface{})["Successful"].(map[string]interface{})["Responses"].([]interface{})
	successfulServerIDs := listSystemsResponse.(map[string]interface{})["Successful"].(map[string]interface{})["ServerIds"].([]interface{})

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

func peripheralServers() map[int64]SystemInfo {
	var minion1ForServer1 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010000-minion-1",
	}
	var minion2ForServer1 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010000-minion-2",
	}
	var peripheralServer1 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010000",
		fqdn: "localhost:8002",
		minions: map[int64]SystemInfo{
			minion1ForServer1.id: minion1ForServer1,
			minion2ForServer1.id: minion2ForServer1,
		},
		port: 8002,
	}

	var minion1ForServer2 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010001-minion-1",
	}
	var minion2ForServer2 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010001-minion-2",
	}
	var peripheralServer2 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010001",
		fqdn: "localhost:8003",
		minions: map[int64]SystemInfo{
			minion1ForServer2.id: minion1ForServer2,
			minion2ForServer2.id: minion2ForServer2,
		},
		port: 8003,
	}

	return map[int64]SystemInfo{
		peripheralServer1.id: peripheralServer1,
		peripheralServer2.id: peripheralServer2,
	}
}
