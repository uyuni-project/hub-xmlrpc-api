package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

//InMemoryHubSessionRepository implements HubSessionRepository
type InMemoryHubSessionRepository struct {
	session *sync.Map
}

func NewInMemoryHubSessionRepository(syncMap *sync.Map) *InMemoryHubSessionRepository {
	return &InMemoryHubSessionRepository{syncMap}
}

func (r *InMemoryHubSessionRepository) SaveHubSession(hubSession *gateway.HubSession) {
	r.session.Store(hubSession.HubSessionKey, hubSession)
}

func (s *InMemoryHubSessionRepository) RetrieveHubSession(hubSessionKey string) *gateway.HubSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession)
	}
	return nil
}

func (s *InMemoryHubSessionRepository) RemoveHubSession(hubSessionKey string) {
	s.session.Delete(hubSessionKey)
}

//InMemoryServerSessionRepository implements ServerSessionRepository
type InMemoryServerSessionRepository struct {
	session *sync.Map
}

func NewInMemoryServerSessionRepository(syncMap *sync.Map) *InMemoryServerSessionRepository {
	return &InMemoryServerSessionRepository{syncMap}
}

func (s *InMemoryServerSessionRepository) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*gateway.ServerSession) {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		for serverID, serverSession := range serverSessions {
			hubSession.(*gateway.HubSession).ServerSessions[serverID] = serverSession
		}
	}
}

func (s *InMemoryServerSessionRepository) RetrieveServerSessions(hubSessionKey string) map[int64]*gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession).ServerSessions
	}
	return make(map[int64]*gateway.ServerSession)
}

func (s *InMemoryServerSessionRepository) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*gateway.HubSession).ServerSessions[serverID]; ok {
			return serverSession
		}
	}
	return nil
}
