package session

import (
	"log"
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/server"
)

type Session struct {
	sessions *sync.Map
	client   server.Client
}

type HubSessionInfo struct {
	username, password string
	loginMode          int
	serverSessionKeys  *sync.Map
}

type ServerSessionInfo struct {
	URL, sessionKey string
}

func (s *Session) SetHubSessionKey(hubSessionKey string, username, password string, loginMode int) {
	s.sessions.Store(hubSessionKey, newHubSessionInfo(username, password, loginMode))
}

func (s *Session) GetLoginMode(hubSessionKey string) int {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSessionInfo.(*HubSessionInfo).loginMode
	}
	return 0
}

func (s *Session) GetUsernameAndPassword(hubSessionKey string) (string, string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSessionInfo.(*HubSessionInfo).username, hubSessionInfo.(*HubSessionInfo).password
	}
	return "", ""
}

func (s *Session) SetServerSessionInfo(hubSessionKey string, serverID int64, serverURL, serverSessionKey string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		hubSessionInfo.(*HubSessionInfo).serverSessionKeys.Store(serverID, &ServerSessionInfo{serverURL, serverSessionKey})
	}
}

func (s *Session) GetServerSessionInfoByServerID(hubSessionKey string, serverID int64) (string, string) {
	if hubSessionInfo, ok := s.sessions.Load(hubSessionKey); ok {
		if serverSessionInfo, ok := hubSessionInfo.(*HubSessionInfo).serverSessionKeys.Load(serverID); ok {
			return serverSessionInfo.(*ServerSessionInfo).URL, serverSessionInfo.(*ServerSessionInfo).sessionKey
		}
	}
	return "", ""
}

func (s *Session) removeHubSessionKey(hubSessionKey string) {
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

func (s *Session) IsHubSessionValid(hubSessionKey string) bool {
	isValid, err := s.client.ExecuteCallToHub("auth.isSessionKeyValid", []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		s.removeHubSessionKey(hubSessionKey)
		return false
	}
	return isValid.(bool)
}

// New returns a new Session struct
func NewSession(client server.Client) *Session {
	var syncMap sync.Map
	return &Session{sessions: &syncMap, client: client}
}
