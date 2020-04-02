package gateway

import "log"

type HubProxy interface {
	DelegateToHub(path string, args []interface{}) (interface{}, error)
}

type HubDelegator struct {
	client         Client
	hubAPIEndpoint string
}

func NewHubDelegator(client Client, hubAPIEndpoint string) *HubDelegator {
	return &HubDelegator{client, hubAPIEndpoint}
}

func (d *HubDelegator) DelegateToHub(path string, args []interface{}) (interface{}, error) {
	response, err := d.client.ExecuteCall(d.hubAPIEndpoint, path, args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return nil, err
	}
	return response, nil
}
