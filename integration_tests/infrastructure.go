package integration_tests

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/rpc"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/parser"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func initInfrastructure(peripheralServersByID map[int64]SystemInfo, port int64, username, password string) {
	initHub(peripheralServersByID, port, username, password)
	initPeripheralServers(peripheralServersByID)
}

func initHub(peripheralServersByID map[int64]SystemInfo, port int64, username, password string) {
	sessionKey := "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4"
	hub := new(UyuniServer)
	hub.mockLogin = func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
		log.Println("Hub -> auth.login", args.Username)
		if args.Username == username && args.Password == password {
			reply.Data = sessionKey
		} else {
			return controller.FaultInvalidCredentials
		}
		return nil
	}
	hub.mockListSystems = func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfoResponse }) error {
		log.Println("Hub -> System.ListSystems", args.SessionKey)
		if args.SessionKey == sessionKey {
			peripheralServers := make([]SystemInfoResponse, 0, len(peripheralServersByID))
			for _, peripheralServer := range peripheralServersByID {
				peripheralServers = append(peripheralServers, SystemInfoResponse{peripheralServer.id, peripheralServer.name})
			}
			reply.Data = peripheralServers
		}
		return nil
	}
	hub.mockListUserSystems = func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfoResponse }) error {
		log.Println("Hub -> System.ListUserSystems", args.Username)
		if args.SessionKey == sessionKey && args.Username == username {
			peripheralServers := make([]SystemInfoResponse, 0, len(peripheralServersByID))
			for _, peripheralServer := range peripheralServersByID {
				peripheralServers = append(peripheralServers, SystemInfoResponse{peripheralServer.id, peripheralServer.name})
			}
			reply.Data = peripheralServers
		}
		return nil
	}
	hub.mockListFqdns = func(r *http.Request, args *struct {
		SessionKey string
		ServerId   int64
	}, reply *struct{ Data []string }) error {
		log.Println("Hub -> System.ListFqdns", args.ServerId)
		if args.SessionKey == sessionKey {
			reply.Data = []string{peripheralServersByID[args.ServerId].fqdn}
		}
		return nil
	}
	initServer(port, hub)
}

func initPeripheralServers(peripheralServersByID map[int64]SystemInfo) {
	for serverID, peripheralServer := range peripheralServersByID {
		serverIDstr := strconv.FormatInt(serverID, 10)
		minions := make([]SystemInfoResponse, 0, len(peripheralServer.minions))
		for _, minion := range peripheralServer.minions {
			minions = append(minions, SystemInfoResponse{minion.id, minion.name})
		}
		server := new(UyuniServer)
		sessionKey := "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd" + serverIDstr
		server.mockLogin = func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
			log.Println("Server"+serverIDstr+" -> auth.login", args.Username)
			reply.Data = sessionKey
			return nil
		}
		server.mockListSystems = func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfoResponse }) error {
			log.Println("Server"+serverIDstr+" -> System.ListSystems", args.SessionKey)
			if args.SessionKey == sessionKey {
				reply.Data = minions
			}
			return nil
		}
		server.mockListUserSystems = func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfoResponse }) error {
			log.Println("Server"+serverIDstr+" -> System.ListUserSystems", args.Username)
			if args.SessionKey == sessionKey && args.Username == "admin" {
				reply.Data = minions
			}
			return nil
		}
		initServer(peripheralServer.port, server)
	}
}

func initServer(port int64, uyuniServer *UyuniServer) {
	go func() {
		rpcServer := rpc.NewServer()
		var codec = xmlrpc.NewCodec()

		codec.RegisterMapping("auth.login", "UyuniServer.Login", parser.LoginRequestParser)
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
	}()
}
