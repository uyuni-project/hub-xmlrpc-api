package gateway

import "log"

type HubProxy interface {
	ProxyCallToHub(call string, args []interface{}) (interface{}, error)
}

type hubProxy struct {
	uyuniHubCallExecutor UyuniHubCallExecutor
}

func NewHubProxy(uyuniHubCallExecutor UyuniHubCallExecutor) *hubProxy {
	return &hubProxy{uyuniHubCallExecutor}
}

func (p *hubProxy) ProxyCallToHub(call string, args []interface{}) (interface{}, error) {
	response, err := p.uyuniHubCallExecutor.ExecuteCall(call, args)
	if err != nil {
		log.Printf("Error ocurred when delegating call to Hub: %v", err)
		return nil, err
	}
	return response, nil
}
