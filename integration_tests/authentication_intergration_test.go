package integration_tests

import (
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

func Test_LoginWithManualMode(t *testing.T) {

	//config.InitializeConfig()
	const errorMessage = "Either the password or username is incorrect"
	tt := []struct {
		name          string
		username      string
		password      string
		expectedError string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "Invalid credentials", username: "unknown", password: "unknown", expectedError: errorMessage},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			client := client.NewClient(10, 10)

			//login
			credentials := []interface{}{tc.username, tc.password}
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.login", credentials)

			// if it's no error then defer the call for closing body
			if err != nil {
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("Unexpected Error message: `%v` doesn't contain `%v`", err, tc.expectedError)
					return
				}
			} else {
				hubSessionKey := loginResponse.(string)
				if len(strings.TrimSpace(hubSessionKey)) == 0 {
					t.Fatalf("Unexpected Result: empty session key was returned.")
				}

			}
		})
	}

}

func Test_LoginWithAutoconnectMode(t *testing.T) {

	//config.InitializeConfig()
	const errorMessage = "Either the password or username is incorrect"
	tt := []struct {
		name          string
		username      string
		password      string
		expectedError string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "Invalid credentials", username: "unknown", password: "unknown", expectedError: errorMessage},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := client.NewClient(10, 10)
			//login
			credentials := []interface{}{tc.username, tc.password}
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAutoconnectMode", credentials)
			if err != nil {
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("Unexpected Error message: `%v` doesn't contain `%v`", err, tc.expectedError)
					return
				}
			} else {
				loginResponseMap := loginResponse.(map[string]interface{})
				hubSessionKey := loginResponseMap["SessionKey"].(string)
				if len(strings.TrimSpace(hubSessionKey)) == 0 {
					t.Fatalf("Unexpected Result: empty session key was returned.")
				}
				validatePeripheralServersResponse(t, loginResponse)
			}
		})
	}

}
func Test_LoginWithAuthRelayMode(t *testing.T) {

	const errorMessage = "Either the password or username is incorrect"
	tt := []struct {
		name          string
		username      string
		password      string
		expectedError string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "Invalid credentials", username: "unknown", password: "unknown", expectedError: errorMessage},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := client.NewClient(10, 10)

			//login
			credentials := []interface{}{tc.username, tc.password}
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAuthRelayMode", credentials)

			if err != nil {
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("Unexpected Error message: `%v` doesn't contain `%v`", err, tc.expectedError)
					return
				}
			} else {
				hubSessionKey := loginResponse.(string)
				if len(strings.TrimSpace(hubSessionKey)) == 0 {
					t.Fatalf("Unexpected Result: empty session key was returned.")
				}

			}
		})
	}

}

func Test_AttachToServers(t *testing.T) {
	const errorMessage = "Either the password or username is incorrect"
	tt := []struct {
		name          string
		username      string
		password      string
		expectedError string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
		{name: "Invalid credentials", username: "unknown", password: "unknown", expectedError: errorMessage},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			client := client.NewClient(10, 10)

			//login
			credentials := []interface{}{tc.username, tc.password}
			loginResponse, err := client.ExecuteCall(gatewayServerURL, "hub.loginWithAuthRelayMode", credentials)

			// if it's no error then defer the call for closing body
			if err != nil {
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("Unexpected Error message: `%v` doesn't contain `%v`", err, tc.expectedError)
					return
				}
			} else {
				hubSessionKey := loginResponse.(string)
				if len(strings.TrimSpace(hubSessionKey)) == 0 {
					t.Fatalf("Unexpected Result: empty session key was returned.")
				}
				loggedInServerIDs := make([]int64, 0, len(peripheralServers))
				for _, serverID := range peripheralServers {
					loggedInServerIDs = append(loggedInServerIDs, serverID.id)
				}
				//Call attachToServers method
				attachToServerResponse, err := client.ExecuteCall(gatewayServerURL, "hub.attachToServers", []interface{}{hubSessionKey, loggedInServerIDs})
				if err != nil {
					t.Fatalf("Unexpected Error message: `%v`", err)
					return

				}
				validatePeripheralServersResponse(t, attachToServerResponse)
			}
		})
	}

}

func validatePeripheralServersResponse(t *testing.T, authResponse interface{}) {
	loginResponseMap := authResponse.(map[string]interface{})
	succeededServers := loginResponseMap["Successful"].(map[string]interface{})
	failedServers := loginResponseMap["Failed"].(map[string]interface{})

	if len(succeededServers) != len(peripheralServers) {
		t.Fatalf("Unexpected Result: Some servers failed unexpectedly")
	}
	if len(failedServers["ServerIds"].([]interface{})) != 0 {
		t.Fatalf("Unexpected Result: there should not be failed severs")
	}

	serverIDs := succeededServers["ServerIds"].([]interface{})
	serverSessionKeys := succeededServers["Responses"].([]interface{})
	for _, sid := range serverIDs {
		if sid.(int64) == 0 {
			t.Fatalf("Unexpected Result: empty session key was returned.")
		}
	}
	for _, skey := range serverSessionKeys {
		if len(strings.TrimSpace(skey.(string))) == 0 {
			t.Fatalf("Unexpected Result: session key was empty.")
		}
	}
}
