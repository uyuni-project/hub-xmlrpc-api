package parser

import (
	"log"
	"reflect"

	"github.com/uyuni-project/hub-xmlrpc-api/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
)

var (
	ListRequestParser      = parseToListRequest
	StructParser           = parseToStruct
	UnicastRequestParser   = parseToUnicastRequest
	MulticastRequestParser = parseToMulitcastRequest
)

func parseToListRequest(request *codec.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.ListRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultInvalidParams
	}
	*parsedArgs = controller.ListRequest{request.MethodName, request.Params}
	return nil
}

func parseToStruct(request *codec.ServerRequest, output interface{}) error {
	val := reflect.ValueOf(output).Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultInvalidParams
	}

	args := request.Params
	if val.NumField() < len(args) {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultWrongArgumentsNumber
	}

	for i, arg := range args {
		field := val.Field(i)
		if field.Type() != reflect.ValueOf(arg).Type() {
			log.Printf("Error ocurred when parsing arguments")
			return codec.FaultInvalidParams
		}
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

func parseToUnicastRequest(request *codec.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.UnicastRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultInvalidParams
	}

	args := request.Params
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return codec.FaultInvalidParams
	}

	serverID, ok := args[1].(int64)
	if !ok {
		log.Printf("Error ocurred when parsing serverID argument")
		return codec.FaultInvalidParams
	}

	rest := args[2:len(args)]
	serverArgs := make([]interface{}, len(rest))
	for i, list := range rest {
		serverArgs[i] = list.(interface{})
	}

	*parsedArgs = controller.UnicastRequest{request.MethodName, hubSessionKey, serverID, serverArgs}
	return nil
}

func parseToMulitcastRequest(request *codec.ServerRequest, output interface{}) error {
	parsedRequest, ok := output.(*controller.MulticastRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultInvalidParams
	}

	args := request.Params
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return codec.FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return codec.FaultInvalidParams
	}

	serverIDs, ok := args[1].([]interface{})
	if !ok {
		log.Printf("Error ocurred when parsing serverIDs argument")
		return codec.FaultInvalidParams
	}

	argsByServer, err := resolveArgsByServer(serverIDs, args[2:len(args)])
	if err != nil {
		return err
	}

	*parsedRequest = controller.MulticastRequest{request.MethodName, hubSessionKey, argsByServer}
	return nil
}

func resolveArgsByServer(serverIDs []interface{}, allServerArgs []interface{}) (map[int64][]interface{}, error) {
	result := make(map[int64][]interface{})
	for i, serverID := range serverIDs {

		parsedServerID, ok := serverID.(int64)
		if !ok {
			log.Printf("Error ocurred when parsing serverIDs argument")
			return nil, codec.FaultInvalidParams
		}

		args := make([]interface{}, 0, len(allServerArgs)+1)

		for _, serverArgs := range allServerArgs {
			parsedServerArgs, ok := serverArgs.([]interface{})
			if !ok {
				log.Printf("Error ocurred when parsing server arguments")
				return nil, codec.FaultInvalidParams
			}
			args = append(args, parsedServerArgs[i])
		}
		result[parsedServerID] = args
	}
	return result, nil
}

func areAllArgumentsOfSameLength(allArrays [][]interface{}) bool {
	//TODO:
	//if !areAllArgumentsOfSameLength(serverArgs) {
	//	return FaultInvalidParams
	//}
	if len(allArrays) <= 1 {
		return true
	}
	lengthToCompare := len(allArrays[0])
	for _, array := range allArrays {
		if lengthToCompare != len(array) {
			return false
		}
	}
	return true
}
