package service

type Client interface {
	ExecuteCall(url string, methodName string, args []interface{}) (response interface{}, err error)
}
