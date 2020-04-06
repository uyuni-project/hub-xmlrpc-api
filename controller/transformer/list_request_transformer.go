package transformer

import (
	"log"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
)

var ListRequestTransformer = transformToListRequest

func transformToListRequest(request *codec.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.ListRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultInvalidParams
	}
	*parsedArgs = controller.ListRequest{request.MethodName, request.Params}
	return nil
}
