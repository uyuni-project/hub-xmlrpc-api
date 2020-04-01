package gateway

type Client interface {
	ExecuteCall(url string, call string, args []interface{}) (response interface{}, err error)
}
