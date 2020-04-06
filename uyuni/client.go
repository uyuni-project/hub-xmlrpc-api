package uyuni

type Client interface {
	ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error)
}
