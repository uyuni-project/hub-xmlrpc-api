package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
)

type SystemInfo struct {
	Id   int64  `xmlrpc:"id"`
	Name string `xmlrpc:"name"`
}

type Auth struct{}
type System struct{}

var System_1 = SystemInfo{
	Id:   1000010000,
	Name: "server-1",
}
var System_2 = SystemInfo{
	Id:   1000010001,
	Name: "server-2",
}
var Systems = []SystemInfo{
	System_1,
	System_2,
}
var sessionkey = "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4"

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	log.Println("Hub -> auth.login", args.Username)
	if args.Username == "admin" && args.Password == "admin" {
		reply.Data = sessionkey
	} else {
		return server.FaultInvalidCredntials
	}
	return nil
}

func (h *Auth) IsSessionKeyValid(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data bool }) error {
	log.Println("Hub -> auth.IsSessionKeyValid", args.SessionKey)
	if args.SessionKey == sessionkey {
		reply.Data = true
	} else {
		return server.Fault{Code: -1, String: "Session id:" + sessionkey + "is not valid."}
	}
	return nil
}

func (h *System) ListUserSystems(r *http.Request, args *struct{ Hubkey, UserLogin string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Hub -> System.ListUserSystems", args.UserLogin)
	if args.Hubkey == sessionkey && args.UserLogin == "admin" {
		reply.Data = Systems
	}
	return nil
}

func (h *System) ListSystems(r *http.Request, args *struct{ Hubkey string }, reply *struct{ Data []SystemInfo }) error {
	log.Println("Hub -> System.ListSystems", args.Hubkey)
	if args.Hubkey == sessionkey {
		reply.Data = Systems
	}
	return nil
}

func (h *System) ListFqdns(r *http.Request, args *struct {
	Hubkey   string
	ServerId int64
}, reply *struct{ Data []string }) error {
	log.Println("Hub -> System.ListFqdns", args.ServerId)
	if args.Hubkey == sessionkey {
		if args.ServerId == 1000010000 {
			reply.Data = []string{"localhost:8002"}
		} else {
			reply.Data = []string{"localhost:8003"}
		}
	}
	return nil
}

func main() {
	RPC := rpc.NewServer()
	var codec = server.NewCodec()
	codec.RegisterDefaultParser(new(server.StructParser))

	codec.RegisterMethod("auth.isSessionKeyValid")
	codec.RegisterMethod("auth.login")
	codec.RegisterMethod("system.listSystems")
	codec.RegisterMethod("system.listUserSystems")
	codec.RegisterMethod("system.listFqdns")

	RPC.RegisterCodec(codec, "text/xml")
	RPC.RegisterService(new(Auth), "auth")
	RPC.RegisterService(new(System), "system")

	//codec.RegisterDefaultMethod("DefaultService.DefaultMethod", new(server.ListParser))

	http.Handle("/hub/rpc/api", RPC)
	log.Println("Starting XML-RPC server on localhost:8001/hub/rpc/api")
	log.Fatal(http.ListenAndServe(":8001", nil))

}
