package gateway

type mockHubSessionRepository struct {
	mockSaveHubSession     func(hubSession *HubSession)
	mockRetrieveHubSession func(hubSessionKey string) *HubSession
	mockRemoveHubSession   func(hubSessionKey string)
}

func (m *mockHubSessionRepository) SaveHubSession(hubSession *HubSession) {
	m.mockSaveHubSession(hubSession)
}
func (m *mockHubSessionRepository) RetrieveHubSession(hubSessionKey string) *HubSession {
	return m.mockRetrieveHubSession(hubSessionKey)
}
func (m *mockHubSessionRepository) RemoveHubSession(hubSessionKey string) {
	m.mockRemoveHubSession(hubSessionKey)
}

type mockServerSessionRepository struct {
	mockSaveServerSessions              func(hubSessionKey string, serverSessions map[int64]*ServerSession)
	mockRetrieveServerSessionByServerID func(hubSessionKey string, serverID int64) *ServerSession
	mockRetrieveServerSessions          func(hubSessionKey string) map[int64]*ServerSession
}

func (m *mockServerSessionRepository) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*ServerSession) {
	m.mockSaveServerSessions(hubSessionKey, serverSessions)
}
func (m *mockServerSessionRepository) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession {
	return m.mockRetrieveServerSessionByServerID(hubSessionKey, serverID)
}
func (m *mockServerSessionRepository) RetrieveServerSessions(hubSessionKey string) map[int64]*ServerSession {
	return m.mockRetrieveServerSessions(hubSessionKey)
}

type mockUyuniAuthenticator struct {
	mockLogin  func(endpoint, username, password string) (string, error)
	mockLogout func(endpoint, sessionKey string) error
}

func (m *mockUyuniAuthenticator) Login(endpoint, username, password string) (string, error) {
	return m.mockLogin(endpoint, username, password)
}

func (m *mockUyuniAuthenticator) Logout(endpoint, sessionKey string) error {
	return m.mockLogout(endpoint, sessionKey)
}

type mockUyuniTopologyInfoRetriever struct {
	mockListServerIDs              func(endpoint, sessionKey string) ([]int64, error)
	mockRetrieveUserServerIDs      func(endpoint, sessionKey, username string) ([]int64, error)
	mockRetrieveServerAPIEndpoints func(endpoint, sessionKey string, serverIDs []int64) (*RetrieveServerAPIEndpointsResponse, error)
}

func (m *mockUyuniTopologyInfoRetriever) ListServerIDs(endpoint, sessionKey string) ([]int64, error) {
	return m.mockListServerIDs(endpoint, sessionKey)
}

func (m *mockUyuniTopologyInfoRetriever) RetrieveUserServerIDs(endpoint, sessionKey, username string) ([]int64, error) {
	return m.mockRetrieveUserServerIDs(endpoint, sessionKey, username)
}

func (m *mockUyuniTopologyInfoRetriever) RetrieveServerAPIEndpoints(endpoint, sessionKey string, serverIDs []int64) (*RetrieveServerAPIEndpointsResponse, error) {
	return m.mockRetrieveServerAPIEndpoints(endpoint, sessionKey, serverIDs)
}

type mockUyuniCallExecutor struct {
	mockExecuteCall func(endpoint string, call string, args []interface{}) (response interface{}, err error)
}

func (m *mockUyuniCallExecutor) ExecuteCall(endpoint string, call string, args []interface{}) (interface{}, error) {
	return m.mockExecuteCall(endpoint, call, args)
}

type mockServerAuthenticator struct {
	mockAttachToServers func(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error)
}

func (m *mockServerAuthenticator) AttachToServers(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error) {
	return m.mockAttachToServers(hubSessionKey, serverIDs, credentialsByServer)
}
