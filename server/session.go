package server

const (
	LOGIN_MANUAL_MODE      = iota // 0
	LOGIN_RELAY_MODE              // 1
	LOGIN_AUTOCONNECT_MODE        // 2
)

type Session interface {
	SetHubSessionKey(hubSessionKey string, username, password string, loginMode int)
	GetLoginMode(hubSessionKey string) int
	GetUsernameAndPassword(hubSessionKey string) (string, string)
	SetServerSessionInfo(hubSessionKey string, serverID int64, serverURL, serverSessionKey string)
	GetServerSessionInfoByServerID(hubSessionKey string, serverID int64) (string, string)
	IsHubSessionValid(hubSessionKey string) bool
}
