package server

const (
	LOGIN_MANUAL_MODE      = iota // 0
	LOGIN_RELAY_MODE              // 1
	LOGIN_AUTOCONNECT_MODE        // 2
)

type ServerSession struct {
	url, sessionKey string
}

type HubSession struct {
	username, password string
	loginMode          int
}

type Session interface {
	SaveHubSession(hubSessionKey string, hubSession *HubSession)
	RetrieveHubSession(hubSessionKey string) *HubSession
	SaveServerSession(hubSessionKey string, serverID int64, serverSessionInfo *ServerSession)
	RetrieveServerSessionByServerID(hubSessionKey string, serverID int64) *ServerSession
	RemoveHubSession(hubSessionKey string)
}
