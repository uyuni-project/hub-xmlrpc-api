package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/service"
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
	hubSession        *service.HubSession
	serverSessionKeys *sync.Map
}

func (s *Session) SaveHubSession(hubSessionKey string, hubSession *service.HubSession) {
	s.sessions.Store(hubSessionKey, newStoredHubSession(hubSession))
}

func (s *Session) RetrieveHubSession(hubSessionKey string) *service.HubSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		return hubSession.(*storedHubSession).hubSession
	}
	return nil
}

func (s *Session) SaveServerSession(hubSessionKey string, serverID int64, serverSession *service.ServerSession) {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		hubSession.(*storedHubSession).serverSessionKeys.Store(serverID, serverSession)
	}
}

func (s *Session) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *service.ServerSession {
	if hubSession, ok := s.sessions.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*storedHubSession).serverSessionKeys.Load(serverID); ok {
			return serverSession.(*service.ServerSession)
		}
	}
	return nil
}

func (s *Session) RemoveHubSession(hubSessionKey string) {
	s.sessions.Delete(hubSessionKey)
}

func newStoredHubSession(hubSession *service.HubSession) *storedHubSession {
	var syncMap sync.Map
	return &storedHubSession{
		hubSession:        hubSession,
		serverSessionKeys: &syncMap,
	}
}
