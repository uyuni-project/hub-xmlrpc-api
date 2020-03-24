package codec

import (
	"reflect"

	"github.com/kolo/xmlrpc"
)

func encodeResponseToXML(response interface{}) (string, error) {
	result := "<methodResponse>"
	params, err := encodeResponseParametersToXML(response)
	if err != nil {
		return "", err
	}
	result += params
	result += "</methodResponse>"
	return result, nil
}

func encodeResponseParametersToXML(response interface{}) (string, error) {
	result := "<params>"
	val := reflect.ValueOf(response).Elem()

	for i := 0; i < val.NumField(); i++ {
		result += "<param>"

		xmlByte, err := xmlrpc.Marshal(val.Field(i).Interface())
		if err != nil {
			return "", err
		}
		result += string(xmlByte)
		result += "</param>"
	}
	result += "</params>"
	return result, nil
}

func encodeFaultToXML(fault Fault) (string, error) {
	result := "<methodResponse><fault>"
	xmlByte, err := xmlrpc.Marshal(fault)
	if err != nil {
		return "", err
	}
	result += string(xmlByte)
	result += "</fault></methodResponse>"
	return result, nil
}
