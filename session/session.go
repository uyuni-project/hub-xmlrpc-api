package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/server"
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
	hubSession        *server.HubSession
	serverSessionKeys *sync.Map
}

func (s *Session) SaveHubSession(hubSessionKey string, hubSession *server.HubSession) {
	s.sessions.Store(hubSessionKey, newStoredHubSession(hubSession))
}

func (s *Session) RetrieveHubSession(hubSessionKey string) *server.HubSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSession.(*storedHubSession).hubSession
	}
	return nil
}

func (s *Session) SaveServerSession(hubSessionKey string, serverID int64, serverSession *server.ServerSession) {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		hubSession.(*storedHubSession).serverSessionKeys.Store(serverID, serverSession)
	}
}

func (s *Session) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *server.ServerSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*storedHubSession).serverSessionKeys.Load(serverID); ok {
			return serverSession.(*server.ServerSession)
		}
	}
	return nil
}

func (s *Session) RemoveHubSession(hubSessionKey string) {
	if _, ok := s.sessions.Load(hubSessionKey); ok {
		s.sessions.Delete(hubSessionKey)
	}
}

func newStoredHubSession(hubSession *server.HubSession) *storedHubSession {
	var syncMap sync.Map
	return &storedHubSession{
		hubSession:        hubSession,
		serverSessionKeys: &syncMap,
	}
}
