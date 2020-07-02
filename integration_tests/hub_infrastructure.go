package integration_tests

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func initInfrastructure(peripheralServersByID map[int64]SystemInfo, port int64, username, password string) {
	initUyuniServer(peripheralServersByID, port, username, password, "hub")
	for _, peripheralServer := range peripheralServersByID {
		initUyuniServer(peripheralServer.minions, peripheralServer.port, username, password, peripheralServer.name)
	}
}

func initUyuniServer(minionsByID map[int64]SystemInfo, port int64, username, password, serverName string) {
	sessionKey := "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4" + serverName
	initServer(port, &UyuniServer{serverName, username, password, sessionKey, minionsByID})
}

func initServer(port int64, uyuniServer *UyuniServer) {
	go func() {
		rpcServer := rpc.NewServer()
		var codec = xmlrpc.NewCodec()

		codec.RegisterMapping("auth.login", "UyuniServer.Login", parser.LoginRequestParser)
		codec.RegisterMapping("auth.logout", "UyuniServer.Logout", parser.LoginRequestParser)
		codec.RegisterMapping("system.listSystems", "UyuniServer.ListSystems", parser.LoginRequestParser)
		codec.RegisterMapping("system.listUserSystems", "UyuniServer.ListUserSystems", parser.LoginRequestParser)
		codec.RegisterMapping("system.listFqdns", "UyuniServer.ListFqdns", parser.LoginRequestParser)

		rpcServer.RegisterCodec(codec, "text/xml")
		rpcServer.RegisterService(uyuniServer, "")

		mux := http.NewServeMux()
		mux.HandleFunc("/rpc/api", func(w http.ResponseWriter, r *http.Request) { rpcServer.ServeHTTP(w, r) })

		log.Printf("Starting XML-RPC server on localhost:%v/rpc/api", port)

		server := http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		}
		log.Fatal(server.ListenAndServe())
		defer server.Close()
	}()
}

type UyuniServer struct {
	serverName, username, password, sessionKey string
	minionsByID                                map[int64]SystemInfo
}

type SystemInfo struct {
	id      int64
	name    string
	fqdn    string
	port    int64
	minions map[int64]SystemInfo
}

type SystemInfoResponse struct {
	Id   int64  `xmlrpc:"id"`
	Name string `xmlrpc:"name"`
}

func (h *UyuniServer) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	log.Println(h.serverName+" -> auth.login", args.Username)
	if args.Username == h.username && args.Password == h.password {
		reply.Data = h.sessionKey
	} else {
		return controller.FaultInvalidCredentials
	}
	return nil
}

func (h *UyuniServer) Logout(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data string }) error {
	log.Println(h.serverName + " -> auth.logout")
	return nil
}

func (u *UyuniServer) ListUserSystems(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfoResponse }) error {
	log.Println(u.serverName+" -> System.ListUserSystems", args.Username)
	if args.SessionKey == u.sessionKey && args.Username == u.username {
		minions := make([]SystemInfoResponse, 0, len(u.minionsByID))
		for _, minion := range u.minionsByID {
			minions = append(minions, SystemInfoResponse{minion.id, minion.name})
		}
		reply.Data = minions
	}
	return nil
}

func (u *UyuniServer) ListSystems(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfoResponse }) error {
	log.Println(u.serverName+" -> System.ListSystems", args.SessionKey)
	if args.SessionKey == u.sessionKey {
		minions := make([]SystemInfoResponse, 0, len(u.minionsByID))
		for _, minion := range u.minionsByID {
			minions = append(minions, SystemInfoResponse{minion.id, minion.name})
		}
		reply.Data = minions
	}
	return nil
}

func (u *UyuniServer) ListFqdns(r *http.Request, args *struct {
	SessionKey string
	ServerId   int64
}, reply *struct{ Data []string }) error {
	log.Println(u.serverName+" -> System.ListFqdns", args.ServerId)
	if args.SessionKey == u.sessionKey {
		reply.Data = []string{u.minionsByID[args.ServerId].fqdn}
	}
	return nil
}
