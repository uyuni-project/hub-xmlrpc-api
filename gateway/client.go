package gateway

type Client interface {
	ExecuteCall(url string, path string, args []interface{}) (response interface{}, err error)
}
