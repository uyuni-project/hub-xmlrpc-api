package gateway

type UyuniAuthenticator interface {
	Login(endpoint, username, password string) (string, error)
	Logout(endpoint, sessionKey string) error
}

type RetrieveServerAPIEndpointsResponse struct {
	SuccessfulResponses map[int64]string
	FailedResponses     map[int64]string
}

type UyuniTopologyInfoRetriever interface {
	ListServerIDs(endpoint, sessionKey string) ([]int64, error)
	RetrieveUserServerIDs(endpoint, sessionKey, username string) ([]int64, error)
	RetrieveServerAPIEndpoints(endpoint, sessionKey string, serverIDs []int64) (*RetrieveServerAPIEndpointsResponse, error)
}

type UyuniCallExecutor interface {
	ExecuteCall(endpoint, call string, args []interface{}) (interface{}, error)
}
