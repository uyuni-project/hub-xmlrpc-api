package gateway

import "log"

type HubProxy interface {
	DelegateToHub(path string, args []interface{}) (interface{}, error)
}

type HubDelegator struct {
	client        Client
	hubSumaAPIURL string
}

func NewHubDelegator(client Client, hubSumaAPIURL string) *HubDelegator {
	return &HubDelegator{client: client, hubSumaAPIURL: hubSumaAPIURL}
}

func (d *HubDelegator) DelegateToHub(path string, args []interface{}) (interface{}, error) {
	response, err := d.client.ExecuteCall(d.hubSumaAPIURL, path, args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return nil, err
	}
	return response, nil
}
