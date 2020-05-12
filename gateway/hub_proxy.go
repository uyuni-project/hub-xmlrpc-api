package gateway

import "log"

type HubProxy interface {
	ProxyCallToHub(call string, args []interface{}) (interface{}, error)
}

type hubProxy struct {
	hubAPIEndpoint    string
	uyuniCallExecutor UyuniCallExecutor
}

func NewHubProxy(hubAPIEndpoint string, uyuniCallExecutor UyuniCallExecutor) *hubProxy {
	return &hubProxy{hubAPIEndpoint, uyuniCallExecutor}
}

func (p *hubProxy) ProxyCallToHub(call string, args []interface{}) (interface{}, error) {
	response, err := p.uyuniCallExecutor.ExecuteCall(p.hubAPIEndpoint, call, args)
	if err != nil {
		log.Printf("Error ocurred when delegating call to Hub: %v", err)
		return nil, err
	}
	return response, nil
}
