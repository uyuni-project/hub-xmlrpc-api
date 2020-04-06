package transformer

import (
	"log"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

var MulticastRequestTransformer = transformToMulitcastRequest

func transformToMulitcastRequest(request *codec.ServerRequest, output interface{}) error {
	parsedRequest, ok := output.(*gateway.MulticastRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultInvalidParams
	}

	args := request.Params
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return controller.FaultInvalidParams
	}

	serverIDs, err := resolveServerIDs(args[1])
	if err != nil {
		return err
	}

	var argsByServer map[int64][]interface{}
	if len(args) > 2 {
		argsByServer, err = resolveArgsByServer(serverIDs, args[2:len(args)])
		if err != nil {
			return err
		}
	}

	method := removeNamespace(request.MethodName)

	*parsedRequest = gateway.MulticastRequest{method, hubSessionKey, serverIDs, argsByServer}
	return nil
}

func resolveServerIDs(args interface{}) ([]int64, error) {
	serverIDs, ok := args.([]interface{})
	if !ok {
		log.Printf("Error ocurred when parsing serverIDs argument")
		return nil, controller.FaultInvalidParams
	}

	parsedServerIDs := make([]int64, 0, len(serverIDs))
	for _, serverID := range serverIDs {
		parsedServerID, ok := serverID.(int64)
		if !ok {
			log.Printf("Error ocurred when parsing serverIDs argument")
			return nil, controller.FaultInvalidParams
		}
		parsedServerIDs = append(parsedServerIDs, parsedServerID)
	}
	return parsedServerIDs, nil
}

func resolveArgsByServer(serverIDs []int64, allServerArgs []interface{}) (map[int64][]interface{}, error) {
	result := make(map[int64][]interface{})
	for i, serverID := range serverIDs {
		args := make([]interface{}, 0, len(allServerArgs)+1)

		for _, serverArgs := range allServerArgs {
			parsedServerArgs, ok := serverArgs.([]interface{})
			if !ok {
				log.Printf("Error ocurred when parsing server arguments")
				return nil, controller.FaultInvalidParams
			}
			args = append(args, parsedServerArgs[i])
		}
		result[serverID] = args
	}
	return result, nil
}

func removeNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
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
