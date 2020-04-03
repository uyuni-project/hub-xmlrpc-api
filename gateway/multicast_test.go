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

func Test_appendServerSessionKeyToServerArgs(t *testing.T) {
	tt := []struct {
		name                 string
		argsByServer         map[int64][]interface{}
		serverSessions       map[int64]*ServerSession
		expectedArgsByServer []serverCall
		expectedErr          string
	}{
		{
			name: "appendServerSessionKeyToServerArgs success",
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			serverSessions: map[int64]*ServerSession{
				1: &ServerSession{1, "1-serverEndpoint", "1-sessionKey", "hubSessionKey"},
				2: &ServerSession{2, "2-serverEndpoint", "2-sessionKey", "hubSessionKey"},
			},
			expectedArgsByServer: []serverCall{
				serverCall{1, "1-serverEndpoint", []interface{}{"1-sessionKey", "arg1_Server1"}},
				serverCall{2, "2-serverEndpoint", []interface{}{"2-sessionKey", "arg1_Server2"}},
			},
		},
		{
			name: "appendServerSessionKeyToServerArgs serverSessions_not_found",
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"arg1_Server1"},
				2: []interface{}{"arg1_Server2"},
			},
			serverSessions: make(map[int64]*ServerSession),
			expectedErr:    "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			mockSession := new(mockSession)
			multicaster := NewMulticaster(new(mockClient), mockSession)

			serverArgs, err := multicaster.appendServerSessionKeyToServerArgs(tc.serverSessions, tc.argsByServer)

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

	tt := []struct {
		name                      string
		argsByServer              map[int64][]interface{}
		mockRetrieveHubSession    func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession
		mockExecuteCall           func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error)
		expectedMulticastResponse *MulticastResponse
		expectedErr               string
	}{
		{
			name:                   "Multicast all_calls_successful",
			argsByServer:           map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return "success_call", nil
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{1: "success_call", 2: "success_call"}, map[int64]interface{}{}},
		},
		{
			name:                   "Multicast all_calls_failed",
			argsByServer:           map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveHubSession: mockRetrieveHubSessionFound,
			mockExecuteCall: func(serverEndpoint string, call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedMulticastResponse: &MulticastResponse{map[int64]interface{}{}, map[int64]interface{}{1: "call_error", 2: "call_error"}},
		},
		{
			name:         "Multicast auth_error invalid_hub_session_key",
			argsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveHubSession: func(argsByServer map[int64][]interface{}) func(hubSessionKey string) *HubSession {
				return func(hubSessionKey string) *HubSession {
					return nil
				}
			},
			expectedErr: "Authentication error: provided session key is invalid",
		},
		{
			name:                   "Multicast serverSessions_not_found",
			argsByServer:           map[int64][]interface{}{1: []interface{}{"arg1_Server1"}, 2: []interface{}{"arg1_Server2"}},
			mockRetrieveHubSession: mockRetrieveHubSessionFoundWithEmptyServerSessions,
			expectedErr:            "Authentication error: provided session key is invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockSession := new(mockSession)
			mockSession.mockRetrieveHubSession = tc.mockRetrieveHubSession(tc.argsByServer)

			mockClient := new(mockClient)
			mockClient.mockExecuteCall = tc.mockExecuteCall

			multicaster := NewMulticaster(mockClient, mockSession)

			multicastResponse, err := multicaster.Multicast("hubSessionKey", "call", tc.argsByServer)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(multicastResponse, tc.expectedMulticastResponse) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedMulticastResponse)
			}
		})
	}
}
