package server

type Client interface {
	ExecuteCall(url string, method string, args []interface{}) (reply interface{}, err error)
}
