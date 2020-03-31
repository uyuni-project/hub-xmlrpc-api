package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type Session struct {
	sessions *sync.Map
}

// NewSession returns a new Session struct
func NewSession() *Session {
	var syncMap sync.Map
	return &Session{sessions: &syncMap}
}

type storedHubSession struct {
	hubSession        *gateway.HubSession
	serverSessionKeys *sync.Map
}

func (s *Session) SaveHubSession(hubSessionKey string, hubSession *gateway.HubSession) {
	s.sessions.Store(hubSessionKey, newStoredHubSession(hubSession))
}

func (s *Session) RetrieveHubSession(hubSessionKey string) *gateway.HubSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSession.(*storedHubSession).hubSession
	}
	return nil
}

func (s *Session) SaveServerSession(hubSessionKey string, serverID int64, serverSession *gateway.ServerSession) {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		hubSession.(*storedHubSession).serverSessionKeys.Store(serverID, serverSession)
	}
}

func (s *Session) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *gateway.ServerSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*storedHubSession).serverSessionKeys.Load(serverID); ok {
			return serverSession.(*gateway.ServerSession)
		}
	}
	return nil
}

func (s *Session) RemoveHubSession(hubSessionKey string) {
	s.sessions.Delete(hubSessionKey)
}

func newStoredHubSession(hubSession *gateway.HubSession) *storedHubSession {
	var syncMap sync.Map
	return &storedHubSession{
		hubSession:        hubSession,
		serverSessionKeys: &syncMap,
	}
}
