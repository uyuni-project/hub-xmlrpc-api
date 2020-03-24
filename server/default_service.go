package server

import (
	"log"
	"net/http"
)

type DefaultService struct {
	client        Client
	hubSumaAPIURL string
}

func NewDefaultService(client Client, hubSumaAPIURL string) *DefaultService {
	return &DefaultService{client: client, hubSumaAPIURL: hubSumaAPIURL}
}

func (d *DefaultService) DefaultMethod(r *http.Request, args *ListArgs, reply *struct{ Data interface{} }) error {
	response, err := d.client.ExecuteCall(d.hubSumaAPIURL, args.Method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
	}
	reply.Data = response
	return nil
}
