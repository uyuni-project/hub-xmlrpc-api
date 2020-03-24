package server

type UnicastArgs struct {
	Method        string
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

type MulticastArgs struct {
	Method        string
	HubSessionKey string
	ServerIDs     []int64
	ServerArgs    [][]interface{}
}

type ListArgs struct {
	Method string
	Args   []interface{}
}

type LoginArgs struct {
	Username string
	Password string
}
