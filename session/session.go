package session

type ApiSession struct {
	sessions map[string]*HubSessionInfo
}

type HubSessionInfo struct {
	username, password string
	serverSessionKeys  map[int64]*struct{ sessionKey, URL string }
}

func (s *ApiSession) SetHubSessionKey(hubSessionKey string, username, password string) {
	s.sessions[hubSessionKey] = NewHubSessionInfo(username, password)
}

//TODO: remove when Abids PR is merged
func (s *ApiSession) IsHubSessionValid(hubSessionKey string) bool {
	_, isValid := s.sessions[hubSessionKey]
	return isValid
}

func (s *ApiSession) GetUsernameAndPassword(hubSessionKey string) (string, string) {
	hubSessionInfo := s.sessions[hubSessionKey]
	return hubSessionInfo.username, hubSessionInfo.password
}

func (s *ApiSession) SetServerSessionInfo(hubSessionKey string, serverID int64, serverURL, serverSessionKey string) {
	s.sessions[hubSessionKey].serverSessionKeys[serverID] = &struct{ sessionKey, URL string }{serverURL, serverSessionKey}
}

func (s *ApiSession) GetServerSessionInfoByServerID(hubSessionKey string, serverID int64) (string, string) {
	serverSessionInfo := s.sessions[hubSessionKey].serverSessionKeys[serverID]
	return serverSessionInfo.URL, serverSessionInfo.sessionKey
}

// New returns a new HubSession struct
func NewHubSessionInfo(username, password string) *HubSessionInfo {
	return &HubSessionInfo{
		username:          username,
		password:          password,
		serverSessionKeys: make(map[int64]*struct{ sessionKey, URL string }),
	}
}

// New returns a new apiSession struct
func New() *ApiSession {
	return &ApiSession{
		sessions: make(map[string]*HubSessionInfo),
	}
}
