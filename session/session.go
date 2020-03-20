package session

import (
	"log"
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

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
	s.sessions.Store(hubSessionKey, newHubSessionInfo(username, password, loginMode))
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

func (s *ApiSession) removeHubSessionKey(hubSessionKey string) {
	if _, ok := s.sessions.Load(hubSessionKey); ok {
		s.sessions.Delete(hubSessionKey)
	}
}

// New returns a new HubSession struct
func newHubSessionInfo(username, password string, loginMode int) *HubSessionInfo {
	var syncMap sync.Map
	return &HubSessionInfo{
		username:          username,
		password:          password,
		loginMode:         loginMode,
		serverSessionKeys: &syncMap,
	}
}

func (s *ApiSession) IsHubSessionValid(hubSessionKey string, client *client.Client) bool {
	isValid, err := client.ExecuteXMLRPCCallToHub("auth.isSessionKeyValid", []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		s.removeHubSessionKey(hubSessionKey)
		return false
	}
	return isValid.(bool)
}

// New returns a new ApiSession struct
func NewApiSession() *ApiSession {
	var syncMap sync.Map
	return &ApiSession{sessions: &syncMap}
}
