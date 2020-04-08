package gateway

type UyuniServerAuthenticator interface {
	Login(endpoint, username, password string) (string, error)
	Logout(endpoint, sessionKey string) error
}

type UyuniServerCallExecutor interface {
	ExecuteCall(endpoint, call string, args []interface{}) (interface{}, error)
}
