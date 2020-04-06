package codec

type Transformer func(request *ServerRequest, output interface{}) error
