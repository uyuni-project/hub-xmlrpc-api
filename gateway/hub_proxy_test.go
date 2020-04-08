package gateway

import (
	"errors"
	"reflect"
	"testing"
)

type mockUyuniHubCallExecutor struct {
	mockExecuteCall func(call string, args []interface{}) (response interface{}, err error)
}

func (m *mockUyuniHubCallExecutor) ExecuteCall(call string, args []interface{}) (interface{}, error) {
	return m.mockExecuteCall(call, args)
}
func Test_ProxyCallToHub(t *testing.T) {
	tt := []struct {
		name             string
		args             []interface{}
		mockExecuteCall  func(call string, args []interface{}) (response interface{}, err error)
		expectedResponse interface{}
		expectedErr      string
	}{
		{
			name: "ProxyCallToHub call_successful",
			args: []interface{}{"arg1", "arg2"},
			mockExecuteCall: func(call string, args []interface{}) (response interface{}, err error) {
				return "success_response", nil
			},
			expectedResponse: "success_response",
		},
		{
			name: "ProxyCallToHub call_error",
			args: []interface{}{"arg1", "arg2"},
			mockExecuteCall: func(call string, args []interface{}) (response interface{}, err error) {
				return nil, errors.New("call_error")
			},
			expectedErr: "call_error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockUyuniHubCallExecutor := new(mockUyuniHubCallExecutor)
			mockUyuniHubCallExecutor.mockExecuteCall = tc.mockExecuteCall

			hubProxy := NewHubProxy(mockUyuniHubCallExecutor)

			response, err := hubProxy.ProxyCallToHub("call", tc.args)

			if err != nil && tc.expectedErr != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
			if err == nil && !reflect.DeepEqual(response, tc.expectedResponse) {
				t.Fatalf("expected and actual don't match, Expected was: %v", tc.expectedResponse)
			}
		})
	}
}
