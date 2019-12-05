package server

import (
	"log"
	"net/http"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/config"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/session"
	"github.com/gorilla/rpc"
)

var conf = config.New()
var apiSession = session.New()

type DefaultService struct{}

type DefaultCallArgs struct {
	HubArgs []interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallArgs, reply *struct{ Data interface{} }) error {
	method, _ := NewCodec().NewRequest(r).Method()
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, method, args.HubArgs)
	if err != nil {
		log.Println("Call error: %v", err)
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

func InitServer() {
	xmlrpcCodec := NewCodec()
	xmlrpcCodec.RegisterMethod("hub.login")
	xmlrpcCodec.RegisterMethod("hub.attachToServers")
	xmlrpcCodec.RegisterMethod("hub.listServerIds")
	xmlrpcCodec.RegisterDefaultMethodForNamespace("multicast", "Multicast.DefaultMethod")
	xmlrpcCodec.RegisterDefaultMethod("DefaultService.DefaultMethod")

	RPC := rpc.NewServer()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(Hub), "hub")
	RPC.RegisterService(new(DefaultService), "")
	RPC.RegisterService(new(Multicast), "")

	http.Handle("/hub/rpc/api", RPC)

	log.Println("Starting XML-RPC server on localhost:8000/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8000", nil))
}