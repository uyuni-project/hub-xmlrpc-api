package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
)

type System struct{}
type Auth struct{}

var sessionkey = "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd6"

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	log.Println("Server2 -> auth.login", args.Username)
	reply.Data = sessionkey
	return nil
}

func main() {
	RPC := rpc.NewServer()
	var codec = server.NewCodec()
	codec.RegisterDefaultParser(new(server.StructParser))

	codec.RegisterMethod("auth.login")

	RPC.RegisterCodec(codec, "text/xml")
	RPC.RegisterService(new(Auth), "auth")

	http.Handle("/rpc/api", RPC)
	log.Println("Starting XML-RPC server on localhost:8003/rpc/api")
	log.Fatal(http.ListenAndServe(":8003", nil))

}
