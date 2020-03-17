package server

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

var conf config.Config

var apiSession = session.New()

type ListArgs struct{ Args []interface{} }

type DefaultService struct{}

func (h *DefaultService) DefaultMethod(r *http.Request, args *ListArgs, reply *struct{ Data interface{} }) error {
	method, _ := NewCodec().NewRequest(r).Method()
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
	}
	reply.Data = response
	return nil
}

func executeXMLRPCCall(url string, method string, args []interface{}) (reply interface{}, err error) {
	client, err := client.GetClientWithTimeout(url, conf.ConnectTimeout, conf.ReadWriteTimeout)
	if err != nil {
		return
	}
	defer client.Close()
	err = client.Call(method, args, &reply)
	return reply, err
}

func InitConfig() {
	conf = config.InitializeConfig()
}

func InitServer() {
	RPC := rpc.NewServer()

	xmlrpcCodec := initXMLRPCCodec()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(Hub), "hub")
	RPC.RegisterService(new(DefaultService), "")
	RPC.RegisterService(new(MulticastService), "")
	RPC.RegisterService(new(Unicast), "")

	http.Handle("/hub/rpc/api", RPC)

	log.Println("Starting XML-RPC server on localhost:8888/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func initXMLRPCCodec() *Codec {
	var codec = NewCodec()

	codec.RegisterDefaultParser(new(StructParser))
	codec.RegisterMethod("hub.login")
	codec.RegisterMethod("hub.loginWithAutoconnectMode")
	codec.RegisterMethod("hub.loginWithAuthRelayMode")
	codec.RegisterMethod("hub.attachToServers")
	codec.RegisterMethod("hub.listServerIds")
	codec.RegisterDefaultMethodForNamespace("multicast", "MulticastService.DefaultMethod", new(MulticastArgsParser))
	codec.RegisterDefaultMethodForNamespace("unicast", "Unicast.DefaultMethod", new(UnicastArgsParser))
	codec.RegisterDefaultMethod("DefaultService.DefaultMethod", new(ListParser))

	return codec
}
