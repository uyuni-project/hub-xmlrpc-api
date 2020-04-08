package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/transformer"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
	"github.com/uyuni-project/hub-xmlrpc-api/uyuni"
)

func main() {
	initServer()
}

func initServer() {
	rpcServer := rpc.NewServer()

	//init config
	conf := config.NewConfig()

	//init xmlrpc client implementation
	client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)

	//init uyuni adapters
	uyuniServerCallExecutor := uyuni.NewUyuniServerCallExecutor(client)
	uyuniServerAuthenticator := uyuni.NewUyuniServerAuthenticator(uyuniServerCallExecutor)

	uyuniHubAuthenticator := uyuni.NewUyuniHubAuthenticator(uyuniServerAuthenticator, conf.Hub.SUMA_API_URL)
	uyuniHubCallExecutor := uyuni.NewUyuniHubCallExecutor(uyuniServerCallExecutor, conf.Hub.SUMA_API_URL)
	uyuniHubTopoloyInfoRetriever := uyuni.NewUyuniHubTopologyInfoRetriever(uyuniHubCallExecutor)

	//init session storage
	session := session.NewInMemorySession()

	//init gateway
	serverAuthenticator := gateway.NewServerAuthenticator(uyuniServerAuthenticator, uyuniHubTopoloyInfoRetriever, session)
	hubAuthenticator := gateway.NewHubAuthenticator(uyuniHubAuthenticator, serverAuthenticator, uyuniHubTopoloyInfoRetriever, session)

	hubProxy := gateway.NewHubProxy(uyuniHubCallExecutor)
	hubTopologyInfoRetriever := gateway.NewHubTopologyInfoRetriever(uyuniHubTopoloyInfoRetriever)

	multicaster := gateway.NewMulticaster(uyuniServerCallExecutor, session)
	unicaster := gateway.NewUnicaster(uyuniServerCallExecutor, session)

	//init controller
	xmlrpcCodec := initCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")

	rpcServer.RegisterService(controller.NewServerAuthenticationController(serverAuthenticator), "")
	rpcServer.RegisterService(controller.NewHubAuthenticationController(hubAuthenticator), "")
	rpcServer.RegisterService(controller.NewHubProxyController(hubProxy), "")
	rpcServer.RegisterService(controller.NewHubTopologyController(hubTopologyInfoRetriever), "")
	rpcServer.RegisterService(controller.NewMulticastController(multicaster), "")
	rpcServer.RegisterService(controller.NewUnicastController(unicaster), "")

	//init server
	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initCodec() *codec.Codec {
	var codec = codec.NewCodec()

	codec.RegisterMapping("hub.login", "HubAuthenticationController.Login", transformer.LoginRequestTransformer)
	codec.RegisterMapping("hub.loginWithAutoconnectMode", "HubAuthenticationController.LoginWithAutoconnectMode", transformer.LoginRequestTransformer)
	codec.RegisterMapping("hub.loginWithAuthRelayMode", "HubAuthenticationController.LoginWithAuthRelayMode", transformer.LoginRequestTransformer)
	codec.RegisterMapping("hub.attachToServers", "ServerAuthenticationController.AttachToServers", transformer.AttachToServersRequestTransformer)
	codec.RegisterMapping("hub.listServerIds", "HubTopologyController.ListServerIDs", transformer.LoginRequestTransformer)

	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastController.Multicast", transformer.MulticastRequestTransformer)
	codec.RegisterDefaultMethodForNamespace("unicast", "UnicastController.Unicast", transformer.UnicastRequestTransformer)
	codec.RegisterDefaultMethod("HubProxyController.ProxyCallToHub", transformer.ListRequestTransformer)

	return codec
}
