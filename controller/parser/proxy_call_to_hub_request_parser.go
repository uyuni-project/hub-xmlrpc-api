package parser

import (
	"log"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func ProxyCallToHubRequestParser(request *xmlrpc.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.ProxyCallToHubRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultInvalidParams
	}
	*parsedArgs = controller.ProxyCallToHubRequest{request.MethodName, request.Params}
	return nil
}
