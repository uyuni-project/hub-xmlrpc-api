package gateway

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_Unicast(t *testing.T) {
	mockIsHubSessionKeyValidTrue := func(hubSessionKey string) bool {
		return true
	}
	mockRetrieveServerSessionByServerIDFound := func(hubSessionKey string, serverID int64) *ServerSession {
		strServerID := strconv.FormatInt(serverID, 10)
		return &ServerSession{serverID, strServerID + "serverAPIEndpoint", strServerID + "serverSessionkey", hubSessionKey}
	}

	tt := []struct {
		name                                string
		serverID                            int64
		serverArgs                          []interface{}
		mockIsHubSessionKeyValid            func(hubSessionKey string) bool
		mockRetrieveServerSessionByServerID func(hubSessionKey string, serverID int64) *ServerSession
		mockExecuteCall                     func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedResponse                    interface{}
		expectedErr                         string
	}{
		{
			name:                                "Unicast call_successful",
			serverID:                            1,
			serverArgs:                          []interface{}{"arg1", "arg2"},
			mockIsHubSessionKeyValid:            mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessionByServerID: mockRetrieveServerSessionByServerIDFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return "success_response", nil
			},
			expectedResponse: "success_response",
		},
		{
			name:                                "Unicast call_error",
			serverID:                            1,
			serverArgs:                          []interface{}{"arg1", "arg2"},
			mockIsHubSessionKeyValid:            mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessionByServerID: mockRetrieveServerSessionByServerIDFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedErr: "call_error",
		},
		{
			name:       "Unicast auth_error invalid_hub_session_key",
			serverID:   1,
			serverArgs: []interface{}{"arg1", "arg2"},
			mockIsHubSessionKeyValid: func(hubSessionKey string) bool {
				return false
			},
			expectedErr: "Authentication error: provided session key is invalid",
		},
		{
			name:                     "Unicast serverSession_not_found",
			serverID:                 1,
			serverArgs:               []interface{}{"arg1", "arg2"},
			mockIsHubSessionKeyValid: mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessionByServerID: func(hubSessionKey string, serverID int64) *ServerSession {
				return nil
			},
			expectedErr: "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockSession := new(mockSession)
			mockSession.mockRetrieveServerSessionByServerID = tc.mockRetrieveServerSessionByServerID

			mockSessionValidator := new(mockSessionValidator)
			mockSessionValidator.mockIsHubSessionKeyValid = tc.mockIsHubSessionKeyValid

			mockClient := new(mockClient)
			mockClient.mockExecuteCall = tc.mockExecuteCall

			unicastService := NewUnicastService(mockClient, mockSession, mockSessionValidator)

			response, err := unicastService.Unicast("hubSessionKey", "call", tc.serverID, tc.serverArgs)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(response, tc.expectedResponse) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedResponse)
			}
		})
	}
}
