package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"io"
	"reflect"

	xmlrpc "github.com/uyuni-project/xmlrpc-public-methods"
)

type Decoder struct {
	*xmlrpc.Decoder
}

func UnmarshalMethodCall(data []byte) (*ServerRequest, error) {
	dec := &Decoder{&xmlrpc.Decoder{xml.NewDecoder(bytes.NewBuffer(data))}}

	if xmlrpc.CharsetReader != nil {
		dec.CharsetReader = xmlrpc.CharsetReader
	}
	serverRequest := &ServerRequest{}
	//process tokens
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if t, ok := token.(xml.StartElement); ok {
			if t.Name.Local == "methodName" {
				if token, err = dec.Token(); err != nil {
					return nil, err
				}
				serverRequest.MethodName = string([]byte(token.(xml.CharData)))
				params, err := dec.unmarshalParameters()
				if err != nil {
					return nil, err
				}
				serverRequest.Params = params
			}
		} else if t, ok := token.(xml.EndElement); ok {
			if t.Name.Local == "methodCall" {
				break
			}
		}
	}
	// read until end of document
	err := dec.Skip()
	if err != nil && err != io.EOF {
		return nil, err
	}
	return serverRequest, nil
}

func (dec *Decoder) unmarshalParameters() ([]interface{}, error) {
	parameters := reflect.ValueOf([]interface{}{})
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if t, ok := token.(xml.StartElement); ok {
			if t.Name.Local == "value" {
				v := reflect.New(parameters.Type().Elem())
				if err := dec.DecodeValue(v); err != nil {
					return nil, err
				}
				parameters = reflect.Append(parameters, v.Elem())
			}
		} else if t, ok := token.(xml.EndElement); ok {
			if t.Name.Local == "params" {
				break
			}
		}
	}
	return parameters.Interface().([]interface{}), nil
}
