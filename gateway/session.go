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
	serverID                                   int64
	serverURL, serverSessionKey, hubSessionKey string
}

func NewServerSession(serverID int64, serverURL, serverSessionKey, hubSessionKey string) *ServerSession {
	return &ServerSession{serverID, serverURL, serverSessionKey, hubSessionKey}
}

type Session interface {
	SaveHubSession(hubSession *HubSession)
	RetrieveHubSession(hubSessionKey string) *HubSession
	RemoveHubSession(hubSessionKey string)
	SaveServerSessions(hubSessionKey string, serverSessions map[int64]*ServerSession)
	RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession
	RetrieveServerSessions(hubSessionKey string) map[int64]*ServerSession
}
