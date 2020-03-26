package codec

type Parser func(request map[string]interface{}, output interface{}) error
