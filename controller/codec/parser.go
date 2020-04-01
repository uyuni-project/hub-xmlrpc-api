package codec

type Parser func(request *ServerRequest, output interface{}) error
