package client

import (
	"reflect"
	"strings"
	"testing"
)

func TestExecuteCall(t *testing.T) {
	tt := []struct {
		name             string
		connectTimeout   int
		requestTimeout   int
		url              string
		methodName       string
		args             []interface{}
		expectedResponse interface{}
		expectedError    string
	}{
		{name: "ExecuteCall Success",
			connectTimeout: 1,
			requestTimeout: 1,
			url:            "http://localhost:8001/hub/rpc/api",
			methodName:     "auth.login",
			args:           []interface{}{"admin", "admin"},
		},
		{name: "ExecuteCall wrong_parameters Failed",
			connectTimeout: 1,
			requestTimeout: 1,
			url:            "http://localhost:8001/hub/rpc/api",
			methodName:     "methodName",
			args:           []interface{}{"wrong_parameter", "wrong_parameter"},
		},
		{name: "ExecuteCall inexistent_URL Failed",
			connectTimeout: 1,
			requestTimeout: 1,
			url:            "http://localhost:8001/hub/rpc/api",
			methodName:     "methodName",
			args:           []interface{}{"admin", "admin"},
		},
		{name: "ExecuteCall inexistent_methodName Failed",
			connectTimeout: 1,
			requestTimeout: 1,
			url:            "http://unknown_URL",
			methodName:     "unknown_method",
			args:           []interface{}{"admin", "admin"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//init server

			//init client
			client := NewClient(tc.connectTimeout, tc.requestTimeout)

			response, err := client.ExecuteCall(tc.url, tc.methodName, tc.args)

			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(response, tc.expectedResponse) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedResponse)
			}
		})
	}
}
