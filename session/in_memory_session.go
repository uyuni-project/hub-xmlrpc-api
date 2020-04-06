package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type InMemorySession struct {
	session *sync.Map
}

func NewInMemorySession() *InMemorySession {
	var session sync.Map
	return &InMemorySession{&session}
}

func (r *InMemorySession) SaveHubSession(hubSession *gateway.HubSession) {
	r.session.Store(hubSession.HubSessionKey, hubSession)
}

func (s *InMemorySession) RetrieveHubSession(hubSessionKey string) *gateway.HubSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession)
	}
	return nil
}

func (s *InMemorySession) RemoveHubSession(hubSessionKey string) {
	s.session.Delete(hubSessionKey)
}

func (s *InMemorySession) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*gateway.ServerSession) {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		for serverID, serverSession := range serverSessions {
			hubSession.(*gateway.HubSession).ServerSessions[serverID] = serverSession
		}
	}
}

func (s *InMemorySession) RetrieveServerSessions(hubSessionKey string) map[int64]*gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession).ServerSessions
	}
	return make(map[int64]*gateway.ServerSession)
}

func (s *InMemorySession) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*gateway.HubSession).ServerSessions[serverID]; ok {
			return serverSession
		}
	}
	return nil
}
