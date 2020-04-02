package gateway

import "log"

type HubProxy interface {
	ProxyCallToHub(call string, args []interface{}) (interface{}, error)
}

type HubDelegator struct {
	client         Client
	hubAPIEndpoint string
}

func NewHubDelegator(client Client, hubAPIEndpoint string) *HubDelegator {
	return &HubDelegator{client, hubAPIEndpoint}
}

func (d *HubDelegator) ProxyCallToHub(call string, args []interface{}) (interface{}, error) {
	response, err := d.client.ExecuteCall(d.hubAPIEndpoint, call, args)
	if err != nil {
		log.Printf("Error ocurred when delegating call to Hub: %v", err)
		return nil, err
	}
	return response, nil
}
