package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/service"
)

type DefaultController struct {
	service *service.DefaultService
}

func NewDefaultController(service *service.DefaultService) *DefaultController {
	return &DefaultController{service}
}

type ListRequest struct {
	Method string
	Args   []interface{}
}

func (d *DefaultController) DefaultMethod(r *http.Request, args *ListRequest, reply *struct{ Data interface{} }) error {
	response, err := d.service.ProcessDefaultCall(args.Method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return err
	}
	reply.Data = response
	return nil
}
