package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubProxyController struct {
	hubProxy gateway.HubProxy
}

type ListRequest struct {
	Method string
	Args   []interface{}
}

func (d *HubProxyController) ProxyCallToHub(r *http.Request, args *ListRequest, reply *struct{ Data interface{} }) error {
	response, err := d.hubProxy.ProxyCallToHub(args.Method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return err
	}
	reply.Data = response
	return nil
}

func NewHubProxyController(hubProxy gateway.HubProxy) *HubProxyController {
	return &HubProxyController{hubProxy}
}
