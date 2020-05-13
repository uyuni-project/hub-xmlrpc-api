package parser

import (
	"log"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
)

var UnicastRequestParser = parseToUnicastRequest

func parseToUnicastRequest(request *codec.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.UnicastRequest)
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

	serverID, ok := args[1].(int64)
	if !ok {
		log.Printf("Error ocurred when parsing serverID argument")
		return controller.FaultInvalidParams
	}

	rest := args[2:len(args)]
	serverArgs := make([]interface{}, len(rest))
	for i, list := range rest {
		serverArgs[i] = list.(interface{})
	}

	method := removeNamespace(request.MethodName)
	*parsedArgs = controller.UnicastRequest{method, hubSessionKey, serverID, serverArgs}
	return nil
}
