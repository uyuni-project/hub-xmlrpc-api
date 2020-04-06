package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type UnicastController struct {
	unicaster gateway.Unicaster
}

func NewUnicastController(unicaster gateway.Unicaster) *UnicastController {
	return &UnicastController{unicaster}
}

func (u *UnicastController) Unicast(r *http.Request, args *gateway.UnicastRequest, reply *struct{ Data interface{} }) error {
	response, err := u.unicaster.Unicast(args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return err
	}
	reply.Data = response
	return nil
}
