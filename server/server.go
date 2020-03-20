package server

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

var apiSession = session.New()

type ListArgs struct{ Args []interface{} }

type DefaultService struct {
	Client *client.Client
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *ListArgs, reply *struct{ Data interface{} }) error {
	method, _ := NewCodec().NewRequest(r).Method()
	response, err := h.Client.ExecuteXMLRPCCallToHub(method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
	}
	reply.Data = response
	return nil
}

func InitServer() {
	rpcServer := rpc.NewServer()

	client := &client.Client{Conf: config.InitializeConfig()}

	xmlrpcCodec := initXMLRPCCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")
	rpcServer.RegisterService(&Hub{Client: client}, "hub")
	rpcServer.RegisterService(&DefaultService{Client: client}, "")
	rpcServer.RegisterService(&MulticastService{Client: client}, "")
	rpcServer.RegisterService(&Unicast{Client: client}, "")

	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initXMLRPCCodec() *Codec {
	var codec = NewCodec()

	codec.RegisterDefaultParser(parseToStruct)
	codec.RegisterMethod("hub.login")
	codec.RegisterMethod("hub.loginWithAutoconnectMode")
	codec.RegisterMethod("hub.loginWithAuthRelayMode")
	codec.RegisterMethodWithParser("hub.attachToServers", parseToMulitcastArgs)
	codec.RegisterMethod("hub.listServerIds")
	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastService.DefaultMethod", parseToMulitcastArgs)
	codec.RegisterDefaultMethodForNamespace("unicast", "Unicast.DefaultMethod", parseToUnicastArgs)
	codec.RegisterDefaultMethod("DefaultService.DefaultMethod", parseToList)

	return codec
}
