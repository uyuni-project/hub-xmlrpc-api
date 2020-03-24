package codec

type Parser func(method string, args []interface{}, output interface{}) error
