package gateway

type UyuniHubAuthenticator interface {
	Login(username, password string) (string, error)
	Logout(hubSessionKey string) error
}

type UyuniHubTopologyInfoRetriever interface {
	ListServerIDs(hubSessionKey string) ([]int64, error)
	RetrieveUserServerIDs(hubSessionKey, username string) ([]int64, error)
	RetrieveServerAPIEndpoints(hubSessionKey string, serverIDs []int64) (map[int64]string, error)
}

type UyuniHubCallExecutor interface {
	ExecuteCall(call string, args []interface{}) (interface{}, error)
}
