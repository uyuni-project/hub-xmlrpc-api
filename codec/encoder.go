package codec

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/kolo/xmlrpc"
)

func encodeResponseToXML(response interface{}) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("<methodResponse><params><param>")

	val := reflect.ValueOf(response).Elem()
	for i := 0; i < val.NumField(); i++ {
		xmlByte, err := xmlrpc.Marshal(val.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		b.WriteString(string(xmlByte))
	}
	b.WriteString("</param></params></methodResponse>")
	return b.Bytes(), nil
}

func encodeFaultErrorToXML(fault Fault) ([]byte, error) {
	var b bytes.Buffer
	xmlByte, err := xmlrpc.Marshal(fault)
	if err != nil {
		return nil, err
	}
	b.WriteString(fmt.Sprintf("<methodResponse><fault>%s</fault></methodResponse>", string(xmlByte)))
	return b.Bytes(), nil
}
