package service

import (
	"log"
)

type DefaultService struct {
	client        Client
	hubSumaAPIURL string
}

func NewDefaultService(client Client, hubSumaAPIURL string) *DefaultService {
	return &DefaultService{client: client, hubSumaAPIURL: hubSumaAPIURL}
}

func (d *DefaultService) ExecuteDefaultCall(path string, args []interface{}) (interface{}, error) {
	response, err := d.client.ExecuteCall(d.hubSumaAPIURL, path, args)
	if err != nil {
		log.Printf("Call error: %v", err)
		return nil, err
	}
	return response, nil
}
