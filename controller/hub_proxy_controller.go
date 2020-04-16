package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubProxyController struct {
	hubProxy gateway.HubProxy
}

type ProxyCallToHubRequest struct {
	Call string
	Args []interface{}
}

func (d *HubProxyController) ProxyCallToHub(r *http.Request, args *ProxyCallToHubRequest, reply *struct{ Data interface{} }) error {
	response, err := d.hubProxy.ProxyCallToHub(args.Call, args.Args)
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
