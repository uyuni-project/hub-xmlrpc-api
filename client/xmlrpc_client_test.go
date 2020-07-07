package client

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	//Hub API Gateway server URL
	sampleResponse = `
	<?xml version="1.0" encoding="UTF-8"?>
	<methodResponse>
		<params>
			<param>
				<value>
				   <string>hello world</string>
				</value>
			</param>
		</params>
	</methodResponse>
    `
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

func TestExecuteCallWithTimeout(t *testing.T) {
	tt := []struct {
		name             string
		connectTimeout   int
		requestTimeout   int
		sleepTime        time.Duration
		expectedResponse interface{}
		expectedError    string
	}{
		{name: "ExecuteCall Pass",
			connectTimeout:   1,
			requestTimeout:   1,
			sleepTime:        0,
			expectedResponse: "hello world",
		},
		{name: "ExecuteCall Fail",
			connectTimeout: 1,
			requestTimeout: 1,
			sleepTime:      2,
			expectedError:  "i/o timeout",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//init server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("Sleeping...")
				time.Sleep(tc.sleepTime * time.Second)
				io.WriteString(w, sampleResponse)
			}))
			defer ts.Close()

			//init client
			client := NewClient(tc.connectTimeout, tc.requestTimeout)
			response, err := client.ExecuteCall(ts.URL, "test", []interface{}{})

			//We expect error
			if len(tc.expectedError) > 0 {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("expected and actual doesn't match, Actuval was: %v, Expected was: %v", err, tc.expectedError)
				}

			} else {
				//We don't expect error
				if err != nil {
					t.Fatalf("Unexpected error was returned: %v", err.Error())
				}
				// We expect a response
				if !reflect.DeepEqual(response, tc.expectedResponse) {
					t.Fatalf("expected and actual doesn't match, Actual was: %v, Expected was: %v", response, tc.expectedResponse)
				}
			}

		})
	}
}
