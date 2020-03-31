package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type UnicastController struct {
	unicaster gateway.Unicaster
}

type UnicastRequest struct {
	Method        string
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

func NewUnicastController(unicaster gateway.Unicaster) *UnicastController {
	return &UnicastController{unicaster}
}

func (u *UnicastController) Unicast(r *http.Request, args *UnicastRequest, reply *struct{ Data interface{} }) error {
	method := removeUnicastNamespace(args.Method)
	response, err := u.unicaster.Unicast(args.HubSessionKey, method, args.ServerID, args.ServerArgs)
	if err != nil {
		log.Printf("Call error: %v", err)
		return err
	}
	reply.Data = response
	return nil
}

func removeUnicastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
