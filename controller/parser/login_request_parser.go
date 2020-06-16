package parser

import (
	"log"
	"reflect"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

var LoginRequestParser = parseToLoginRequest

func parseToLoginRequest(request *xmlrpc.ServerRequest, output interface{}) error {
	val := reflect.ValueOf(output).Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultInvalidParams
	}

	args := request.Params
	if val.NumField() < len(args) {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultWrongArgumentsNumber
	}

	for i, arg := range args {
		field := val.Field(i)
		if field.Type() != reflect.ValueOf(arg).Type() {
			log.Printf("Error ocurred when parsing arguments")
			return controller.FaultInvalidParams
		}
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}
