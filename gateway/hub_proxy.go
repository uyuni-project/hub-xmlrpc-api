package gateway

import "log"

type HubProxy interface {
	ProxyCallToHub(call string, args []interface{}) (interface{}, error)
}

type hubProxy struct {
	client         Client
	hubAPIEndpoint string
}

func NewHubProxy(client Client, hubAPIEndpoint string) *hubProxy {
	return &hubProxy{client, hubAPIEndpoint}
}

func (d *hubProxy) ProxyCallToHub(call string, args []interface{}) (interface{}, error) {
	response, err := d.client.ExecuteCall(d.hubAPIEndpoint, call, args)
	if err != nil {
		log.Printf("Error ocurred when delegating call to Hub: %v", err)
		return nil, err
	}
	return response, nil
}
