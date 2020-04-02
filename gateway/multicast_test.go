package gateway

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

type mockSession struct {
	mockSaveHubSession                  func(hubSession *HubSession)
	mockRetrieveHubSession              func(hubSessionKey string) *HubSession
	mockRemoveHubSession                func(hubSessionKey string)
	mockSaveServerSessions              func(hubSessionKey string, serverSessions map[int64]*ServerSession)
	mockRetrieveServerSessionByServerID func(hubSessionKey string, serverID int64) *ServerSession
	mockRetrieveServerSessions          func(hubSessionKey string) map[int64]*ServerSession
}

func (m *mockSession) SaveHubSession(hubSession *HubSession) {
	m.mockSaveHubSession(hubSession)
}
func (m *mockSession) RetrieveHubSession(hubSessionKey string) *HubSession {
	return m.mockRetrieveHubSession(hubSessionKey)
}
func (m *mockSession) RemoveHubSession(hubSessionKey string) {
	m.mockRemoveHubSession(hubSessionKey)
}
func (m *mockSession) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*ServerSession) {
	m.mockSaveServerSessions(hubSessionKey, serverSessions)
}
func (m *mockSession) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession {
	return m.mockRetrieveServerSessionByServerID(hubSessionKey, serverID)
}
func (m *mockSession) RetrieveServerSessions(hubSessionKey string) map[int64]*ServerSession {
	return m.mockRetrieveServerSessions(hubSessionKey)
}

type mockClient struct {
	mockExecuteCall func(endpoint string, call string, args []interface{}) (response interface{}, err error)
}

func (m *mockClient) ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error) {
	return m.mockExecuteCall(endpoint, call, args)
}

type mockSessionValidator struct {
	mockIsHubSessionKeyValid func(hubSessionKey string) bool
}

func (m *mockSessionValidator) isHubSessionKeyValid(hubSessionKey string) bool {
	return m.mockIsHubSessionKeyValid(hubSessionKey)
}

func Test_appendServerSessionKeyToServerArgs(t *testing.T) {
	tt := []struct {
		name                       string
		argsByServer               map[int64][]interface{}
		mockRetrieveServerSessions func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession
		expectedArgsByServer       []serverCall
		expectedErr                string
	}{
		{
			name:         "appendServerSessionKeyToServerArgs success",
			argsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveServerSessions: func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession {
				return func(hubSessionKey string) map[int64]*ServerSession {
					result := make(map[int64]*ServerSession)
					for serverID := range argsByServer {
						strServerID := strconv.FormatInt(serverID, 10)
						result[serverID] = &ServerSession{serverID, strServerID + "-serverEndpoint", strServerID + "-sessionKey", hubSessionKey}
					}
					return result
				}
			},
			expectedArgsByServer: []serverCall{
				serverCall{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
				serverCall{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
			},
		},
		{
			name:         "appendServerSessionKeyToServerArgs serverSessions_not_found",
			argsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveServerSessions: func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession {
				return func(hubSessionKey string) map[int64]*ServerSession {
					return make(map[int64]*ServerSession)
				}
			},
			expectedErr: "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			mockSession := new(mockSession)
			mockSession.mockRetrieveServerSessions = tc.mockRetrieveServerSessions(tc.argsByServer)
			multicastService := NewMulticastService(new(mockClient), mockSession, new(mockSessionValidator))

			serverArgs, err := multicastService.appendServerSessionKeyToServerArgs("hubSessionKey", tc.argsByServer)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(serverArgs, tc.expectedArgsByServer) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedArgsByServer)
			}
		})
	}
}

func Test_executeCallOnServers(t *testing.T) {
	tt := []struct {
		name                      string
		serverCalls               []serverCall
		mockExecuteCall           func(serverCalls []serverCall) func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedMulticastResponse *MulticastResponse
	}{
		{
			name: "executeCallOnServers all_calls_successful",
			serverCalls: []serverCall{
				serverCall{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
				serverCall{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
			},
			mockExecuteCall: func(serverCalls []serverCall) func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
					return "success_call", nil
				}
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{1: "success_call", 2: "success_call"}, map[int64]interface{}{}},
		},
		{
			name: "executeCallOnServers first_call_successful_and_the_other_calls_failed",
			serverCalls: []serverCall{
				serverCall{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
				serverCall{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
			},
			mockExecuteCall: func(serverCalls []serverCall) func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
					if serverCalls[0].serverEndpoint == serverEndpoint {
						return "success_call", nil
					}
					return nil, errors.New("call_error")
				}
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{1: "success_call"}, map[int64]interface{}{2: "call_error"}},
		},
		{
			name: "executeCallOnServers all_calls_failed",
			serverCalls: []serverCall{
				serverCall{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
				serverCall{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
			},
			mockExecuteCall: func(serverCalls []serverCall) func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
					return nil, errors.New("call_error")
				}
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{}, map[int64]interface{}{1: "call_error", 2: "call_error"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mockClient)
			mockClient.mockExecuteCall = tc.mockExecuteCall(tc.serverCalls)

			multicastResponse := executeCallOnServers("call", tc.serverCalls, mockClient)

			if !reflect.DeepEqual(multicastResponse, tc.expectedMulticastResponse) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedMulticastResponse)
			}
		})
	}
}

func Test_Multicast(t *testing.T) {
	mockIsHubSessionKeyValidTrue := func(hubSessionKey string) bool {
		return true
	}

	mockRetrieveServerSessionsFound := func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession {
		return func(hubSessionKey string) map[int64]*ServerSession {
			result := make(map[int64]*ServerSession)
			for serverID := range argsByServer {
				strServerID := strconv.FormatInt(serverID, 10)
				result[serverID] = &ServerSession{serverID, strServerID + "-serverEndpoint", strServerID + "-sessionKey", hubSessionKey}
			}
			return result
		}
	}

	tt := []struct {
		name                       string
		argsByServer               map[int64][]interface{}
		mockIsHubSessionKeyValid   func(hubSessionKey string) bool
		mockRetrieveServerSessions func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession
		mockExecuteCall            func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedMulticastResponse  *MulticastResponse
		expectedErr                string
	}{
		{
			name:                       "Multicast all_calls_successful",
			argsByServer:               map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockIsHubSessionKeyValid:   mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessions: mockRetrieveServerSessionsFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return "success_call", nil
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{1: "success_call", 2: "success_call"}, map[int64]interface{}{}},
		},
		{
			name:                       "Multicast all_calls_failed",
			argsByServer:               map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockIsHubSessionKeyValid:   mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessions: mockRetrieveServerSessionsFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{}, map[int64]interface{}{1: "call_error", 2: "call_error"}},
		},
		{
			name:         "Multicast auth_error invalid_hub_session_key",
			argsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockIsHubSessionKeyValid: func(hubSessionKey string) bool {
				return false
			},
			mockRetrieveServerSessions: mockRetrieveServerSessionsFound,
			expectedErr:                "Authentication error: provided session key is invalid",
		},
		{
			name:                     "Multicast serverSessions_not_found",
			argsByServer:             map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockIsHubSessionKeyValid: mockIsHubSessionKeyValidTrue,
			mockRetrieveServerSessions: func(argsByServer map[int64][]interface{}) func(hubSessionKey string) map[int64]*ServerSession {
				return func(hubSessionKey string) map[int64]*ServerSession {
					return make(map[int64]*ServerSession)
				}
			},
			expectedErr: "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockSession := new(mockSession)
			mockSession.mockRetrieveServerSessions = tc.mockRetrieveServerSessions(tc.argsByServer)

			mockSessionValidator := new(mockSessionValidator)
			mockSessionValidator.mockIsHubSessionKeyValid = tc.mockIsHubSessionKeyValid

			mockClient := new(mockClient)
			mockClient.mockExecuteCall = tc.mockExecuteCall

			multicastService := NewMulticastService(mockClient, mockSession, mockSessionValidator)

			multicastResponse, err := multicastService.Multicast("hubSessionKey", "call", tc.argsByServer)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(multicastResponse, tc.expectedMulticastResponse) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedMulticastResponse)
			}
		})
	}
}
