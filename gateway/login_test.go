package gateway

import (
	"errors"
	"reflect"
	"testing"
)

func Test_Login(t *testing.T) {
	tt := []struct {
		name                  string
		mockLogin             func(hubSessionKey string) func(endpoint, username, password string) (string, error)
		expectedHubSessionKey string
		expectedErr           string
	}{
		{
			name: "Login success",
			mockLogin: func(hubSessionKey string) func(endpoint, username, password string) (string, error) {
				return func(endpoint, username, password string) (string, error) {
					return hubSessionKey, nil
				}
			},
			expectedHubSessionKey: "hubSessionKey",
		},
		{
			name: "Login login_to_hub_server_failed ",
			mockLogin: func(hubSessionKey string) func(endpoint, username, password string) (string, error) {
				return func(endpoint, username, password string) (string, error) {
					return "", errors.New("login_error")
				}
			},
			expectedErr: "login_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockHubSessionRepository := new(mockHubSessionRepository)
			savedOnSession := false
			mockHubSessionRepository.mockSaveHubSession = func(hubSession *HubSession) { savedOnSession = true }

			mockUyuniAuthenticator := new(mockUyuniAuthenticator)
			mockUyuniAuthenticator.mockLogin = tc.mockLogin(tc.expectedHubSessionKey)
			mockUyuniTopologyInfoRetrtiever := new(mockUyuniTopologyInfoRetriever)
			mockServerAuthenticator := new(mockServerAuthenticator)

			hubLoginer := NewHubLoginer("hub_API_endpoint", mockUyuniAuthenticator, mockServerAuthenticator, mockUyuniTopologyInfoRetrtiever, mockHubSessionRepository)

			hubSessionKey, err := hubLoginer.Login("username", "password")

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(hubSessionKey, tc.expectedHubSessionKey) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedHubSessionKey)
			}
			if err == nil && !savedOnSession {
				t.Fatalf("HubSession was not saved as expected")
			}
		})
	}
}

func Test_LoginWithAutoconnectMode(t *testing.T) {
	mockLoginToHubServerSuccess := func(hubSessionKey string) func(endpoint, username, password string) (string, error) {
		return func(endpoint, username, password string) (string, error) {
			return hubSessionKey, nil
		}
	}
	mockRetrieveUserServerIDsSuccess := func(endpoint, sessionKey, username string) ([]int64, error) {
		return nil, nil
	}
	attachToServersResponse := &MulticastResponse{
		map[int64]ServerSuccessfulResponse{
			1: ServerSuccessfulResponse{1, "1-serverEndpoint", "success_call"},
		},
		map[int64]ServerFailedResponse{
			2: ServerFailedResponse{2, "2-serverEndpoint", "failed_call"},
		},
	}

	tt := []struct {
		name                                     string
		mockLogin                                func(hubSessionKey string) func(endpoint, username, password string) (string, error)
		mockRetrieveUserServerIDs                func(endpoint, sessionKey, username string) ([]int64, error)
		mockAttachToServers                      func(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error)
		hubSessionKey                            string
		attachToServersResponse                  *MulticastResponse
		expectedLoginWithAutoconnectModeResponse *LoginWithAutoconnectModeResponse
		expectedErr                              string
	}{
		{
			name:                      "LoginWithAutoconnectMode success",
			mockLogin:                 mockLoginToHubServerSuccess,
			mockRetrieveUserServerIDs: mockRetrieveUserServerIDsSuccess,
			mockAttachToServers: func(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error) {
				return attachToServersResponse, nil
			},
			hubSessionKey: "hubSessionKey",
			expectedLoginWithAutoconnectModeResponse: &LoginWithAutoconnectModeResponse{
				"hubSessionKey", attachToServersResponse,
			},
		},
		{
			name: "LoginWithAutoconnectMode login_to_hub_server_failed",
			mockLogin: func(hubSessionKey string) func(endpoint, username, password string) (string, error) {
				return func(endpoint, username, password string) (string, error) {
					return "", errors.New("login_error")
				}
			},
			expectedErr: "login_error",
		},
		{
			name:      "LoginWithAutoconnectMode retrieve_the_user_serverIDs_failed ",
			mockLogin: mockLoginToHubServerSuccess,
			mockRetrieveUserServerIDs: func(endpoint, sessionKey, username string) ([]int64, error) {
				return nil, errors.New("retrieve_the_user_serverIDs_error")
			},
			expectedErr: "retrieve_the_user_serverIDs_error",
		},
		{
			name:                      "LoginWithAutoconnectMode attach_peripheral_servers_sessions_to_hub_session_failed ",
			mockLogin:                 mockLoginToHubServerSuccess,
			mockRetrieveUserServerIDs: mockRetrieveUserServerIDsSuccess,
			mockAttachToServers: func(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error) {
				return nil, errors.New("attach_to_servers_error")
			},
			expectedErr: "attach_to_servers_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockHubSessionRepository := new(mockHubSessionRepository)
			savedOnSession := false
			mockHubSessionRepository.mockSaveHubSession = func(hubSession *HubSession) { savedOnSession = true }

			mockUyuniAuthenticator := new(mockUyuniAuthenticator)
			mockUyuniAuthenticator.mockLogin = tc.mockLogin(tc.hubSessionKey)
			mockUyuniTopologyInfoRetrtiever := new(mockUyuniTopologyInfoRetriever)
			mockUyuniTopologyInfoRetrtiever.mockRetrieveUserServerIDs = tc.mockRetrieveUserServerIDs
			mockServerAuthenticator := new(mockServerAuthenticator)
			mockServerAuthenticator.mockAttachToServers = tc.mockAttachToServers

			hubLoginer := NewHubLoginer("hub_API_endpoint", mockUyuniAuthenticator, mockServerAuthenticator, mockUyuniTopologyInfoRetrtiever, mockHubSessionRepository)

			loginWithAutoconnectModeResponse, err := hubLoginer.LoginWithAutoconnectMode("username", "password")

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(loginWithAutoconnectModeResponse, tc.expectedLoginWithAutoconnectModeResponse) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedLoginWithAutoconnectModeResponse)
			}
			if err == nil && !savedOnSession {
				t.Fatalf("HubSession was not saved as expected")
			}
		})
	}
}
