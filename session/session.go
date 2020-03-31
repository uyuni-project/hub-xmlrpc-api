package session

import (
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type Session struct {
	session *sync.Map
}

func NewSession(session *sync.Map) *Session {
	return &Session{session}
}

func (r *Session) SaveHubSession(hubSession *gateway.HubSession) {
	r.session.Store(hubSession.HubSessionKey, hubSession)
}

func (s *Session) RetrieveHubSession(hubSessionKey string) *gateway.HubSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession)
	}
	return nil
}

func (s *Session) RemoveHubSession(hubSessionKey string) {
	s.session.Delete(hubSessionKey)
}

func (s *Session) SaveServerSessions(hubSessionKey string, serverSessions map[int64]*gateway.ServerSession) {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		for serverID, serverSession := range serverSessions {
			hubSession.(*gateway.HubSession).ServerSessions[serverID] = serverSession
		}
	}
}

func (s *Session) RetrieveServerSessions(hubSessionKey string) map[int64]*gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		return hubSession.(*gateway.HubSession).ServerSessions
	}
	return make(map[int64]*gateway.ServerSession)
}

func (s *Session) RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *gateway.ServerSession {
	if hubSession, ok := s.session.Load(hubSessionKey); ok {
		if serverSession, ok := hubSession.(*gateway.HubSession).ServerSessions[serverID]; ok {
			return serverSession
		}
	}
	return nil
}
