package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
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

	xmlrpcCodec := initXMLRPCCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")
	rpcServer.RegisterService(server.NewHubService(client, session, conf.Hub.SUMA_API_URL), "hub")
	rpcServer.RegisterService(server.NewDefaultService(client, conf.Hub.SUMA_API_URL), "")
	rpcServer.RegisterService(server.NewMulticastService(client, session, conf.Hub.SUMA_API_URL), "")
	rpcServer.RegisterService(server.NewUnicastService(client, session, conf.Hub.SUMA_API_URL), "")

	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initXMLRPCCodec() *codec.Codec {
	var codec = codec.NewCodec()

	codec.RegisterDefaultParser(parser.StructParser)
	codec.RegisterMethod("hub.login")
	codec.RegisterMethod("hub.loginWithAutoconnectMode")
	codec.RegisterMethod("hub.loginWithAuthRelayMode")
	codec.RegisterMethodWithParser("hub.attachToServers", parser.MulticastParser)
	codec.RegisterMethod("hub.listServerIds")
	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastService.DefaultMethod", parser.MulticastParser)
	codec.RegisterDefaultMethodForNamespace("unicast", "Unicast.DefaultMethod", parser.UnicastParser)
	codec.RegisterDefaultMethod("DefaultService.DefaultMethod", parser.ListParser)

	return codec
}
