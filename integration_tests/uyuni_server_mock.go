package integration_tests

import (
	"net/http"
)

type UyuniServer struct {
	mockLogin           func(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error
	mockListUserSystems func(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfoResponse }) error
	mockListSystems     func(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfoResponse }) error
	mockListFqdns       func(r *http.Request, args *struct {
		SessionKey string
		ServerId   int64
	}, reply *struct{ Data []string }) error
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
	return h.mockLogin(r, args, reply)
}

func (h *UyuniServer) ListUserSystems(r *http.Request, args *struct{ SessionKey, Username string }, reply *struct{ Data []SystemInfoResponse }) error {
	return h.mockListUserSystems(r, args, reply)
}

func (h *UyuniServer) ListSystems(r *http.Request, args *struct{ SessionKey string }, reply *struct{ Data []SystemInfoResponse }) error {
	return h.mockListSystems(r, args, reply)
}

func (h *UyuniServer) ListFqdns(r *http.Request, args *struct {
	SessionKey string
	ServerId   int64
}, reply *struct{ Data []string }) error {
	return h.mockListFqdns(r, args, reply)
}
