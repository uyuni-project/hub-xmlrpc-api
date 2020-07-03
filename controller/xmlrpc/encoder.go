package xmlrpc

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	xmlrpc "github.com/uyuni-project/xmlrpc-public-methods"
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

func encodeFaultErrorToXML(fault controller.FaultError) ([]byte, error) {
	var b bytes.Buffer
	xmlByte, err := xmlrpc.Marshal(fault)
	if err != nil {
		return nil, err
	}
	b.WriteString(fmt.Sprintf("<methodResponse><fault>%s</fault></methodResponse>", string(xmlByte)))
	return b.Bytes(), nil
}
