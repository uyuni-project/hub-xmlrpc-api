package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/service"
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

	xmlrpcCodec := initCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")
	rpcServer.RegisterService(controller.NewHubController(service.NewHubService(client, session, conf.Hub.SUMA_API_URL)), "")
	rpcServer.RegisterService(controller.NewDefaultController(service.NewDefaultService(client, conf.Hub.SUMA_API_URL)), "")
	rpcServer.RegisterService(controller.NewMulticastController(service.NewMulticastService(client, session, conf.Hub.SUMA_API_URL)), "")
	rpcServer.RegisterService(controller.NewUnicastController(service.NewUnicastService(client, session, conf.Hub.SUMA_API_URL)), "")

	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initCodec() *codec.Codec {
	var codec = codec.NewCodec()

	codec.RegisterDefaultParser(parser.StructParser)

	codec.RegisterMapping("hub.login", "HubController.Login")
	codec.RegisterMapping("hub.loginWithAutoconnectMode", "HubController.LoginWithAutoconnectMode")
	codec.RegisterMapping("hub.loginWithAuthRelayMode", "HubController.LoginWithAuthRelayMode")
	codec.RegisterMapping("hub.listServerIds", "HubController.ListServerIds")
	codec.RegisterMappingWithParser("hub.attachToServers", "HubController.AttachToServers", parser.MulticastParser)

	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastController.DefaultMethod", parser.MulticastParser)
	codec.RegisterDefaultMethodForNamespace("unicast", "UnicastController.DefaultMethod", parser.UnicastParser)
	codec.RegisterDefaultMethod("DefaultController.DefaultMethod", parser.ListParser)

	return codec
}
