package gateway

type HubSession struct {
	HubSessionKey, username, password string
	loginMode                         int
	ServerSessions                    map[int64]*ServerSession
}

func NewHubSession(hubSessionKey, username, password string, loginMode int) *HubSession {
	return &HubSession{hubSessionKey, username, password, loginMode, make(map[int64]*ServerSession)}
}

type ServerSession struct {
	serverID                                           int64
	serverAPIEndpoint, serverSessionKey, hubSessionKey string
}

func NewServerSession(serverID int64, serverEndpoinit, serverSessionKey, hubSessionKey string) *ServerSession {
	return &ServerSession{serverID, serverEndpoinit, serverSessionKey, hubSessionKey}
}

type HubSessionRepository interface {
	SaveHubSession(hubSession *HubSession)
	RetrieveHubSession(hubSessionKey string) *HubSession
	RemoveHubSession(hubSessionKey string)
}

type ServerSessionRepository interface {
	SaveServerSessions(hubSessionKey string, serverSessions map[int64]*ServerSession)
	RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession
	RetrieveServerSessions(hubSessionKey string) map[int64]*ServerSession
}
