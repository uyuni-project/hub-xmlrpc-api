package session

import "sync"

type ApiSession struct {
	sessions *sync.Map
}

const (
	LOGIN_MANUAL_MODE      = iota // 0
	LOGIN_RELAY_MODE              // 1
	LOGIN_AUTOCONNECT_MODE        // 2
)

type HubSessionInfo struct {
	username, password string
	loginMode          int
	serverSessionKeys  *sync.Map
}

type ServerSessionInfo struct {
	URL, sessionKey string
}

func (s *ApiSession) SetHubSessionKey(hubSessionKey string, username, password string, loginMode int) {
	s.sessions.Store(hubSessionKey, NewHubSessionInfo(username, password, loginMode))
}

func (s *ApiSession) RemoveHubSessionKey(hubSessionKey string) {
	if _, ok := s.sessions.Load(hubSessionKey); ok {
		s.sessions.Delete(hubSessionKey)
	}
}

func (s *ApiSession) GetLoginMode(hubSessionKey string) int {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSessionInfo.(*HubSessionInfo).loginMode
	}
	return 0
}

func (s *ApiSession) GetUsernameAndPassword(hubSessionKey string) (string, string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSessionInfo.(*HubSessionInfo).username, hubSessionInfo.(*HubSessionInfo).password
	}
	return "", ""
}

func (s *ApiSession) SetServerSessionInfo(hubSessionKey string, serverID int64, serverURL, serverSessionKey string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		hubSessionInfo.(*HubSessionInfo).serverSessionKeys.Store(serverID, &ServerSessionInfo{serverURL, serverSessionKey})
	}
}

func (s *ApiSession) GetServerSessionInfoByServerID(hubSessionKey string, serverID int64) (string, string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		if serverSessionInfo, ok := hubSessionInfo.(*HubSessionInfo).serverSessionKeys.Load(serverID); ok {
			return serverSessionInfo.(*ServerSessionInfo).URL, serverSessionInfo.(*ServerSessionInfo).sessionKey
		}
	}
	return "", ""
}

// New returns a new HubSession struct
func NewHubSessionInfo(username, password string, loginMode int) *HubSessionInfo {
	var syncMap sync.Map
	return &HubSessionInfo{
		username:          username,
		password:          password,
		loginMode:         loginMode,
		serverSessionKeys: &syncMap,
	}
}

// New returns a new ApiSession struct
func New() *ApiSession {
	var syncMap sync.Map
	return &ApiSession{
		sessions: &syncMap,
	}
}
