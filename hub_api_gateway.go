package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
	"github.com/uyuni-project/hub-xmlrpc-api/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

func main() {
	initServer()
}

func initServer() {
	rpcServer := rpc.NewServer()

	conf := config.InitializeConfig()
	client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
	session := session.NewSession()
	authorizer := gateway.NewAuthorizationService(client, session, conf.Hub.SUMA_API_URL)

	xmlrpcCodec := initCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")

	rpcServer.RegisterService(controller.NewAuthorizerController(gateway.NewAuthorizationService(client, session, conf.Hub.SUMA_API_URL)), "")
	rpcServer.RegisterService(controller.NewHubProxyController(gateway.NewHubDelegator(client, conf.Hub.SUMA_API_URL)), "")
	rpcServer.RegisterService(controller.NewHubController(gateway.NewHubServiceImpl(client, session, conf.Hub.SUMA_API_URL, authorizer)), "")
	rpcServer.RegisterService(controller.NewMulticastController(gateway.NewMulticastService(client, session, authorizer)), "")
	rpcServer.RegisterService(controller.NewUnicastController(gateway.NewUnicastService(client, session, authorizer)), "")

	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initCodec() *codec.Codec {
	var codec = codec.NewCodec()

	codec.RegisterDefaultParser(parser.StructParser)

	codec.RegisterMapping("hub.login", "AuthorizerController.Login")
	codec.RegisterMapping("hub.loginWithAutoconnectMode", "AuthorizerController.LoginWithAutoconnectMode")
	codec.RegisterMapping("hub.loginWithAuthRelayMode", "AuthorizerController.LoginWithAuthRelayMode")
	codec.RegisterMappingWithParser("hub.attachToServers", "AuthorizerController.AttachToServers", parser.MulticastRequestParser)

	codec.RegisterMapping("hub.listServerIds", "HubController.ListServerIDs")

	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastController.Multicast", parser.MulticastRequestParser)
	codec.RegisterDefaultMethodForNamespace("unicast", "UnicastController.Unicast", parser.UnicastRequestParser)
	codec.RegisterDefaultMethod("HubProxyController.DelegateToHub", parser.ListRequestParser)

	return codec
}
