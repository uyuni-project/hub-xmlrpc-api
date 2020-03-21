package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

func main() {
	initServer()
}

func initServer() {
	rpcServer := rpc.NewServer()

	client := client.NewClient(config.InitializeConfig())
	apiSession := session.NewSession(client)

	xmlrpcCodec := initXMLRPCCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")
	rpcServer.RegisterService(server.NewHubService(client, apiSession), "hub")
	rpcServer.RegisterService(server.NewDefaultService(client), "")
	rpcServer.RegisterService(server.NewMulticastService(client, apiSession), "")
	rpcServer.RegisterService(server.NewUnicastService(client, apiSession), "")

	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initXMLRPCCodec() *server.Codec {
	var codec = server.NewCodec()

	codec.RegisterDefaultParser(server.StructParser)
	codec.RegisterMethod("hub.login")
	codec.RegisterMethod("hub.loginWithAutoconnectMode")
	codec.RegisterMethod("hub.loginWithAuthRelayMode")
	codec.RegisterMethodWithParser("hub.attachToServers", server.MulitcastParser)
	codec.RegisterMethod("hub.listServerIds")
	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastService.DefaultMethod", server.MulitcastParser)
	codec.RegisterDefaultMethodForNamespace("unicast", "Unicast.DefaultMethod", server.UnicastParser)
	codec.RegisterDefaultMethod("DefaultService.DefaultMethod", server.ListParser)

	return codec
}
