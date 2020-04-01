package gateway

import (
	"testing"
	"time"
)

var (
	tm              = time.Now()
	mockExecuteCall func(url string, path string, args []interface{}) (response interface{}, err error)

	mockSaveHubSession                  func(hubSession *HubSession)
	mockRetrieveHubSession              func(hubSessionKey string) *HubSession
	mockRemoveHubSession                func(hubSessionKey string)
	mockSaveServerSessions              func(hubSessionKey string, serverSessions map[int64]*ServerSession)
	mockRetrieveServerSessionByServerID func(hubSessionKey string, serverID int64) *ServerSession
	mockRetrieveServerSessions          func(hubSessionKey string) map[int64]*ServerSession

	mockIsHubSessionKeyValid func(hubSessionKey string) bool
)

type MockSession struct{}

func (c *MockSession) SaveHubSession(hubSession *HubSession) {
	mockSaveHubSession(hubSession)
}
func (c *MockSession) RetrieveHubSession(hubSessionKey string) *HubSession {
	return mockRetrieveHubSession(hubSessionKey)
}
func (c *MockSession) RemoveHubSession(hubSessionKey string) {
	mockRemoveHubSession(hubSessionKey)
}
func (c *MockSession) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*ServerSession) {
	mockSaveServerSessions(hubSessionKey, serverSessions)
}
func (c *MockSession) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession {
	return mockRetrieveServerSessionByServerID(hubSessionKey, serverID)
}
func (c *MockSession) RetrieveServerSessions(hubSessionKey string) map[int64]*ServerSession {
	return mockRetrieveServerSessions(hubSessionKey)
}

type MockClient struct{}

func (c *MockClient) ExecuteCall(url string, path string, args []interface{}) (response interface{}, err error) {
	return mockExecuteCall(url, path, args)
}

type MockSessionValidator struct{}

func (a *MockSessionValidator) isHubSessionKeyValid(hubSessionKey string) bool {
	return mockIsHubSessionKeyValid(hubSessionKey)
}

func Test_appendServerSessionKeyToServerArgs(t *testing.T) {

	//success test setup
	type serverArgs struct {
		serverURL, serverSessionKey string
		serverID                    int64
	}
	serverArgsT := []serverArgs{serverArgs{"url", "serverSessionkey", 1}}
	argsByServer := map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}}

	mockRetrieveServerSessions = func(hubSessionKey string) map[int64]*ServerSession {
		result := make(map[int64]*ServerSession)

		for args := range serverArgsT {
			if args.serverSessionKey != "" {

			}
			result[serverID] = &ServerSession{serverID, serverURL, serverSessionKey, hubSessionKey}
		}
	}

	tt := []struct {
		name                       string
		hubSessionKey              string
		call                       string
		argsByServer               map[int64][]interface{}
		mockRetrieveServerSessions func(hubSessionKey string) map[int64]*ServerSession
		expectedErr                string
	}{
		{name: "appendServerSessionKeyToServerArgs Success",
			hubSessionKey: "HubSessionKey",
			call:          "call_to_execute",
			argsByServer:  map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			mockSession := new(MockSession)
			mockRetrieveServerSessions = tc.mockRetrieveServerSessions
			mockClient := new(MockClient)
			mockSessionValidator := new(MockSessionValidator)
			multicastService := NewMulticastService(mockClient, mockSession, mockSessionValidator)

			serverArgs, err := multicastService.appendServerSessionKeyToServerArgs(tc.hubSessionKey, tc.argsByServer)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
		})
	}
}

func Test_Multicast(t *testing.T) {
	tt := []struct {
		name          string
		hubSessionKey string
		call          string
		argsByServer  map[int64][]interface{}
	}{
		{name: "appendServerSessionKeyToServerArgs Success",
			hubSessionKey: "HubSessionKey",
			call:          "call_to_execute",
			argsByServer:  map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			mockSession := new(MockSession)
			mockClient := new(MockClient)
			mockSessionValidator := new(MockSessionValidator)
			multicastService := NewMulticastService(mockClient, mockSession, mockSessionValidator)

			multicastResponse, err := multicastService.Multicast(tc.hubSessionKey, tc.call, tc.argsByServer)

			if err != nil && tc.output != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
		})
	}
}
