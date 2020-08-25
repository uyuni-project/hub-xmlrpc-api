package gateway

import (
	"errors"
	"strings"
	"testing"
)

func Test_Logout(t *testing.T) {
	mockRetrieveHubSessionFound := func(hubSessionKey string) *HubSession {
		return &HubSession{"hubSessionKey", "username", "password", 1, make(map[int64]*ServerSession)}
	}
	tt := []struct {
		name                   string
		mockUyuniServerLogout  func(endpoint, sessionKey string) error
		mockRetrieveHubSession func(hubSessionKey string) *HubSession
		expectedError          string
	}{
		{
			name: "Logout_should_succeed",
			mockUyuniServerLogout: func(endpoint, sessionKey string) error {
				return nil
			},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
		},
		{
			name: "Logout no_session_found_should_fail ",
			mockUyuniServerLogout: func(endpoint, sessionKey string) error {
				return errors.New("logout_error")
			},
			mockRetrieveHubSession: func(hubSessionKey string) *HubSession {
				return nil
			},
			expectedError: "Authentication error: provided session key is invalid",
		},
		{
			name: "Logout error_when_logging_out_from_hub_server_should_fail ",
			mockUyuniServerLogout: func(endpoint, sessionKey string) error {
				return errors.New("logout_error")
			},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
			expectedError:          "logout_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockHubSessionRepository := new(mockHubSessionRepository)
			sessionRemoved := false
			mockHubSessionRepository.mockRemoveHubSession = func(hubSessionKey string) { sessionRemoved = true }
			mockHubSessionRepository.mockRetrieveHubSession = tc.mockRetrieveHubSession

			mockUyuniAuthenticator := new(mockUyuniAuthenticator)
			mockUyuniAuthenticator.mockLogout = tc.mockUyuniServerLogout

			hubLogouter := NewHubLogouter("hub_API_endpoint", mockUyuniAuthenticator, mockHubSessionRepository)

			err := hubLogouter.Logout("hubSessionKey")

			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("expected and actual errors don't match. Expected was:\n%v\nActual is:\n%v:", tc.expectedError, err.Error())
			}
			if err == nil && !sessionRemoved {
				t.Fatalf("HubSession was not removed as expected")
			}
		})
	}
}
