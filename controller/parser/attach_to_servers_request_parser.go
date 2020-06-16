package parser

import (
	"log"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

func AttachToServersRequestParser(request *xmlrpc.ServerRequest, output interface{}) error {
	parsedRequest, ok := output.(*controller.AttachToServersRequest)
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

	var credentialsByServer map[int64]*gateway.Credentials
	if len(args) > 2 {
		credentialsByServer, err = resolveCredentialsByServer(serverIDs, args[2:len(args)])
		if err != nil {
			return err
		}
	}

	*parsedRequest = controller.AttachToServersRequest{hubSessionKey, serverIDs, credentialsByServer}
	return nil
}

func resolveCredentialsByServer(serverIDs []int64, allServerArgs []interface{}) (map[int64]*gateway.Credentials, error) {
	if len(allServerArgs) != 2 {
		log.Printf("Error ocurred when parsing credentials")
		return nil, controller.FaultInvalidParams
	}
	usernames, ok := allServerArgs[0].([]interface{})
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return nil, controller.FaultInvalidParams
	}
	passwords, ok := allServerArgs[1].([]interface{})
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return nil, controller.FaultInvalidParams
	}

	result := make(map[int64]*gateway.Credentials)
	for i, serverID := range serverIDs {
		username, ok := usernames[i].(string)
		if !ok {
			log.Printf("Error ocurred when parsing arguments")
			return nil, controller.FaultInvalidParams
		}
		password, ok := passwords[i].(string)
		if !ok {
			log.Printf("Error ocurred when parsing arguments")
			return nil, controller.FaultInvalidParams
		}
		result[serverID] = &gateway.Credentials{username, password}
	}
	return result, nil
}
