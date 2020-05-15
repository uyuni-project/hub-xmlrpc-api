package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubLogoutController struct {
	hubLogouter gateway.HubLogouter
}

func NewHubLogoutController(hubLogouter gateway.HubLogouter) *HubLogoutController {
	return &HubLogoutController{hubLogouter}
}

type LogoutRequest struct {
	HubSessionKey string
}

func (h *HubLogoutController) Logout(r *http.Request, args *LogoutRequest, reply *struct{ Data string }) error {
	err := h.hubLogouter.Logout(args.HubSessionKey)
	if err != nil {
		log.Printf("Logout error: %v", err)
		return err
	}
	return nil
}
