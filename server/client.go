package server

type Client interface {
	ExecuteCallWithURL(url string, method string, args []interface{}) (reply interface{}, err error)
	ExecuteCallToHub(method string, args []interface{}) (reply interface{}, err error)
}
