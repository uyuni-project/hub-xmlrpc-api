package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
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
	conf := config.InitConfig()

	//init xmlrpc client implementation
	client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)

	//init uyuni adapters
	uyuniCallExecutor := uyuni.NewUyuniCallExecutor(client)
	uyuniAuthenticator := uyuni.NewUyuniAuthenticator(uyuniCallExecutor)
	uyuniTopoloyInfoRetriever := uyuni.NewUyuniTopologyInfoRetriever(uyuniCallExecutor)

	//init session storage
	var syncMap sync.Map
	hubSessionRepository := session.NewInMemoryHubSessionRepository(&syncMap)
	serverSessionRepository := session.NewInMemoryServerSessionRepository(&syncMap)

	//init gateway
	serverAuthenticator := gateway.NewServerAuthenticator(conf.Hub.SUMA_API_URL, uyuniAuthenticator, uyuniTopoloyInfoRetriever, hubSessionRepository, serverSessionRepository)
	hubLoginer := gateway.NewHubLoginer(conf.Hub.SUMA_API_URL, uyuniAuthenticator, serverAuthenticator, uyuniTopoloyInfoRetriever, hubSessionRepository)
	hubLogouter := gateway.NewHubLogouter(conf.Hub.SUMA_API_URL, uyuniAuthenticator, hubSessionRepository)

	hubProxy := gateway.NewHubProxy(conf.Hub.SUMA_API_URL, uyuniCallExecutor)
	hubTopologyInfoRetriever := gateway.NewTopologyInfoRetriever(conf.Hub.SUMA_API_URL, uyuniTopoloyInfoRetriever)

	multicaster := gateway.NewMulticaster(uyuniCallExecutor, hubSessionRepository)
	unicaster := gateway.NewUnicaster(uyuniCallExecutor, serverSessionRepository)

	//init controller
	xmlrpcCodec := initCodec()
	rpcServer.RegisterCodec(xmlrpcCodec, "text/xml")

	rpcServer.RegisterService(controller.NewServerAuthenticationController(serverAuthenticator, transformer.MulticastResponseTransformer), "")
	rpcServer.RegisterService(controller.NewHubLoginController(hubLoginer, transformer.MulticastResponseTransformer), "")
	rpcServer.RegisterService(controller.NewHubLogoutController(hubLogouter), "")
	rpcServer.RegisterService(controller.NewHubProxyController(hubProxy), "")
	rpcServer.RegisterService(controller.NewHubTopologyController(hubTopologyInfoRetriever), "")
	rpcServer.RegisterService(controller.NewMulticastController(multicaster, transformer.MulticastResponseTransformer), "")
	rpcServer.RegisterService(controller.NewUnicastController(unicaster), "")

	//init server
	http.Handle("/hub/rpc/api", rpcServer)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initCodec() *codec.Codec {
	var codec = codec.NewCodec()

	codec.RegisterMapping("hub.login", "HubLoginController.Login", parser.LoginRequestParser)
	codec.RegisterMapping("hub.loginWithAutoconnectMode", "HubLoginController.LoginWithAutoconnectMode", parser.LoginRequestParser)
	codec.RegisterMapping("hub.loginWithAuthRelayMode", "HubLoginController.LoginWithAuthRelayMode", parser.LoginRequestParser)
	codec.RegisterMapping("hub.logout", "HubLogoutController.Logout", parser.LoginRequestParser)
	codec.RegisterMapping("hub.attachToServers", "ServerAuthenticationController.AttachToServers", parser.AttachToServersRequestParser)
	codec.RegisterMapping("hub.listServerIds", "HubTopologyController.ListServerIDs", parser.LoginRequestParser)

	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastController.Multicast", parser.MulticastRequestParser)
	codec.RegisterDefaultMethodForNamespace("unicast", "UnicastController.Unicast", parser.UnicastRequestParser)
	codec.RegisterDefaultMethod("HubProxyController.ProxyCallToHub", parser.ProxyCallToHubRequestParser)

	return codec
}
