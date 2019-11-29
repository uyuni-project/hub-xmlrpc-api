package serverencoder

import (
	"reflect"

	"github.com/kolo/xmlrpc"
)

func encodeResponseToXML(response interface{}) (string, error) {
	buffer := "<methodResponse>"
	params, err := encodeResponseParametersToXML(response)
	buffer += params
	buffer += "</methodResponse>"
	return buffer, err
}

func encodeResponseParametersToXML(response interface{}) (string, error) {
	var err error
	buffer := "<params>"

	val := reflect.ValueOf(response)

	for i := 0; i < val.Elem().NumField(); i++ {
		var xmlByte []byte
		buffer += "<param>"
		xmlByte, err = xmlrpc.Marshal(val.Elem().Field(i).Interface())
		buffer += string(xmlByte)
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer, err
}
