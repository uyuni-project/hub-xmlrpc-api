package gateway

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_executeCallOnServers(t *testing.T) {
	tt := []struct {
		name                      string
		multicastCallRequest      *multicastCallRequest
		expectedMulticastResponse *MulticastResponse
	}{
		{
			name: "executeCallOnServers all_calls_successful",
			multicastCallRequest: &multicastCallRequest{
				func(endpoint string, args []interface{}) (interface{}, error) {
					return "success_call", nil
				},
				[]serverCallInfo{
					serverCallInfo{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
					serverCallInfo{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
				},
			},
			expectedMulticastResponse: &MulticastResponse{
				map[int64]ServerSuccessfulResponse{
					1: ServerSuccessfulResponse{1, "1-serverEndpoint", "success_call"},
					2: ServerSuccessfulResponse{2, "2-serverEndpoint", "success_call"},
				},
				map[int64]ServerFailedResponse{},
			},
		},
		{
			name: "executeCallOnServers first_call_successful_and_the_other_calls_failed",
			multicastCallRequest: &multicastCallRequest{
				func(endpoint string, args []interface{}) (interface{}, error) {
					if endpoint == "1-serverEndpoint" {
						return "success_call", nil
					}
					return nil, errors.New("call_error")
				},
				[]serverCallInfo{
					serverCallInfo{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
					serverCallInfo{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
				},
			},
			expectedMulticastResponse: &MulticastResponse{
				map[int64]ServerSuccessfulResponse{
					1: ServerSuccessfulResponse{1, "1-serverEndpoint", "success_call"},
				},
				map[int64]ServerFailedResponse{
					2: ServerFailedResponse{2, "2-serverEndpoint", "call_error"},
				},
			},
		},
		{
			name: "executeCallOnServers all_calls_failed",
			multicastCallRequest: &multicastCallRequest{
				func(endpoint string, args []interface{}) (interface{}, error) {
					return nil, errors.New("call_error")
				},
				[]serverCallInfo{
					serverCallInfo{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
					serverCallInfo{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
				},
			},
			expectedMulticastResponse: &MulticastResponse{
				map[int64]ServerSuccessfulResponse{},
				map[int64]ServerFailedResponse{
					1: ServerFailedResponse{1, "1-serverEndpoint", "call_error"},
					2: ServerFailedResponse{2, "2-serverEndpoint", "call_error"},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multicastResponse := executeCallOnServers(tc.multicastCallRequest)

			if !reflect.DeepEqual(multicastResponse, tc.expectedMulticastResponse) {
				t.Fatalf("expected and actual don't match. Actual was:  %v. Expected was: %v", multicastResponse, tc.expectedMulticastResponse)
			}
		})
	}
}

func Test_Multicast(t *testing.T) {
	mockRetrieveHubSessionFound := func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession {
		return func(hubSessionKey string) *HubSession {
			serverSessions := make(map[int64]*ServerSession)
			for serverID := range argsByServer {
				strServerID := strconv.FormatInt(serverID, 10)
				serverSessions[serverID] =
					&ServerSession{serverID, strServerID + "-serverEndpoint", strServerID + "-sessionKey", hubSessionKey}
			}
			return &HubSession{"hubSessionKey", "username", "password", 1, serverSessions}
		}
	}
	mockRetrieveHubSessionFoundWithEmptyServerSessions :=
		func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession {
			return func(hubSessionKey string) *HubSession {
				return &HubSession{"hubSessionKey", "username", "password", 1, make(map[int64]*ServerSession)}
			}
		}

	mockRetrieveHubSessionNotFound := func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession {
		return func(hubSessionKey string) *HubSession {
			return nil
		}
	}
	mockExecuteCallSuccessful := func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
		return "success_call", nil
	}
	mockExecuteCallError := func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
		return nil, errors.New("call_error")
	}

	tt := []struct {
		name                      string
		serverIDs                 []int64
		argsByServer              map[int64][]interface{}
		mockRetrieveHubSession    func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession
		mockExecuteCall           func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedMulticastResponse *MulticastResponse
		expectedErr               string
	}{
		{
			name:      "Multicast all_calls_successful",
			serverIDs: []int64{1, 2},
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
			mockExecuteCall:        mockExecuteCallSuccessful,
			expectedMulticastResponse: &MulticastResponse{
				map[int64]ServerSuccessfulResponse{
					1: ServerSuccessfulResponse{1, "1-serverEndpoint", "success_call"},
					2: ServerSuccessfulResponse{2, "2-serverEndpoint", "success_call"},
				},
				map[int64]ServerFailedResponse{},
			},
		},
		{
			name:      "Multicast all_calls_failed",
			serverIDs: []int64{1, 2},
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
			mockExecuteCall:        mockExecuteCallError,
			expectedMulticastResponse: &MulticastResponse{
				map[int64]ServerSuccessfulResponse{},
				map[int64]ServerFailedResponse{
					1: ServerFailedResponse{1, "1-serverEndpoint", "call_error"},
					2: ServerFailedResponse{2, "2-serverEndpoint", "call_error"},
				},
			},
		},
		{
			name:      "Multicast auth_error invalid_hub_session_key",
			serverIDs: []int64{1, 2},
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			mockRetrieveHubSession: mockRetrieveHubSessionNotFound,
			expectedErr:            "Authentication error: provided session key is invalid",
		},
		{
			name:      "Multicast serverSessions_not_found",
			serverIDs: []int64{1, 2},
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			mockRetrieveHubSession: mockRetrieveHubSessionFoundWithEmptyServerSessions,
			expectedErr:            "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockSession := new(mockHubSessionRepository)
			mockSession.mockRetrieveHubSession = tc.mockRetrieveHubSession(tc.argsByServer)

			mockUyuniCallExecutor := new(mockUyuniCallExecutor)
			mockUyuniCallExecutor.mockExecuteCall = tc.mockExecuteCall

			multicaster := NewMulticaster(mockUyuniCallExecutor, mockSession)

			multicastResponse, err := multicaster.Multicast("hubSessionKey", "call", tc.serverIDs, tc.argsByServer)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(multicastResponse, tc.expectedMulticastResponse) {
				t.Fatalf("Expected and actual values don't match, Expected value is: %v", tc.expectedMulticastResponse)
			}
		})
	}
}
