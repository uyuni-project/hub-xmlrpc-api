package uyuni

type uyuniCallExecutor struct {
	client Client
}

func NewUyuniCallExecutor(client Client) *uyuniCallExecutor {
	return &uyuniCallExecutor{client}
}

func (u *uyuniCallExecutor) ExecuteCall(endpoint, call string, args []interface{}) (interface{}, error) {
	response, err := u.client.ExecuteCall(endpoint, call, args)
	if err != nil {
		return "", err
	}
	return response, nil
}
