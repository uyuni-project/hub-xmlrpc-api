package gateway

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

type mockUyuniCallExecutor struct {
	mockExecuteCall func(endpoint string, call string, args []interface{}) (response interface{}, err error)
}

func (m *mockUyuniCallExecutor) ExecuteCall(endpoint string, call string, args []interface{}) (interface{}, error) {
	return m.mockExecuteCall(endpoint, call, args)
}
