package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func Test_LoginRequestParser(t *testing.T) {
	type testLoginRequest struct{ Username, Password string }

	tt := []struct {
		name             string
		serverRequest    *xmlrpc.ServerRequest
		requestToHydrate interface{}
		expectedRequest  testLoginRequest
		expectedError    string
	}{
		{name: "LoginRequestParser should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"username", "password"}},
			requestToHydrate: &testLoginRequest{},
			expectedRequest:  testLoginRequest{Username: "username", Password: "password"}},
		{name: "LoginRequestParser no_struct_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"username", "password"}},
			requestToHydrate: &[]interface{}{},
			expectedRequest:  testLoginRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "LoginRequestParser wrong_number_of_arguments_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"username", "password", "extra_argument"}},
			requestToHydrate: &testLoginRequest{},
			expectedRequest:  testLoginRequest{},
			expectedError:    controller.FaultWrongArgumentsNumber.Message},
		{name: "LoginRequestParser wrong_type_of_arguments_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"username", 123}},
			requestToHydrate: &testLoginRequest{},
			expectedRequest:  testLoginRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := LoginRequestParser(tc.serverRequest, tc.requestToHydrate)
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("expected and actual errors don't match. Expected was:\n%v\nActual is:\n%v:", tc.expectedError, err.Error())
			}
			if err == nil && !reflect.DeepEqual(tc.requestToHydrate, &tc.expectedRequest) {
				t.Fatalf("expected and actual requests don't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedRequest, tc.requestToHydrate)
			}
		})
	}
}

func Test_ProxyCallToHubRequestParser(t *testing.T) {
	tt := []struct {
		name             string
		serverRequest    *xmlrpc.ServerRequest
		requestToHydrate interface{}
		expectedRequest  controller.ProxyCallToHubRequest
		expectedError    string
	}{
		{name: "ProxyCallToHubRequestParser should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			requestToHydrate: &controller.ProxyCallToHubRequest{},
			expectedRequest:  controller.ProxyCallToHubRequest{Call: "method", Args: []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}}},
		{name: "ProxyCallToHubRequestParser nil_parameters_should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"sessionKey", "arg1_Hub", nil}},
			requestToHydrate: &controller.ProxyCallToHubRequest{},
			expectedRequest:  controller.ProxyCallToHubRequest{Call: "method", Args: []interface{}{"sessionKey", "arg1_Hub", nil}}},
		{name: "ProxyCallToHubRequestParser no_ProxyCallToHubRequest_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.ProxyCallToHubRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := ProxyCallToHubRequestParser(tc.serverRequest, tc.requestToHydrate)
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("expected and actual errors don't match. Expected was:\n%v\nActual is:\n%v:", tc.expectedError, err.Error())
			}
			if err == nil && !reflect.DeepEqual(tc.requestToHydrate, &tc.expectedRequest) {
				t.Fatalf("expected and actual requests don't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedRequest, tc.requestToHydrate)
			}
		})
	}
}

func Test_MulticastRequestParser(t *testing.T) {
	tt := []struct {
		name             string
		serverRequest    *xmlrpc.ServerRequest
		requestToHydrate interface{}
		expectedRequest  controller.MulticastRequest
		expectedError    string
	}{
		{name: "MulticastRequestParser single_server_should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1)}, []interface{}{"arg1_Server1"}, []interface{}{"arg2_Server1"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{Call: "method", HubSessionKey: "hubSessionKey", ServerIDs: []int64{1}, ArgsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}}}},
		{name: "MulticastRequestParser multiple_server_should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{Call: "method", HubSessionKey: "hubSessionKey", ServerIDs: []int64{1, 2}, ArgsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}}}},
		{name: "MulticastRequestParser nil_parameters_should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{nil, nil}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{Call: "method", HubSessionKey: "hubSessionKey", ServerIDs: []int64{1, 2}, ArgsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1", nil}, 2: []interface{}{"arg1_Server2", nil}}}},
		{name: "MulticastRequestParser no_serverID_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{"serverSessionKey"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser malformed_hubSessionKey_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{1, []interface{}{int64(1)}, []interface{}{"arg1_Server1"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser malformed_serverID_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{"not_a_number", "not_a_number"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.MulticastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser no_MulticastRequest_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1)}, "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.MulticastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser no_namespace_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedError:    controller.FaultDecode.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := MulticastRequestParser(tc.serverRequest, tc.requestToHydrate)
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("expected and actual errors don't match. Expected was:\n%v\nActual is:\n%v:", tc.expectedError, err.Error())
			}
			if err == nil && !reflect.DeepEqual(tc.requestToHydrate, &tc.expectedRequest) {
				t.Fatalf("expected and actual structs don't match. Expected was:\n%v\nActual is:\n%v:", &tc.expectedRequest, tc.requestToHydrate)
			}
		})
	}
}

func Test_UnicastRequestParser(t *testing.T) {
	tt := []struct {
		name             string
		serverRequest    *xmlrpc.ServerRequest
		requestToHydrate interface{}
		expectedRequest  controller.UnicastRequest
		expectedError    string
	}{
		{name: "UnicastRequestParser should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.UnicastRequest{Call: "method", HubSessionKey: "sessionKey", ServerID: int64(1), Args: []interface{}{"arg1_Server1", "arg2_Server1"}}},
		{name: "UnicastRequestParser nil_parameters_should_succeed",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", int64(1), "arg1", nil}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.UnicastRequest{Call: "method", HubSessionKey: "sessionKey", ServerID: int64(1), Args: []interface{}{"arg1", nil}}},
		{name: "UnicastRequestParser wrong_number_of_arguments Failed",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.UnicastRequest{},
			expectedError:    controller.FaultWrongArgumentsNumber.Message},
		{name: "UnicastRequestParser malformed_hubSessionKey_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{123, int64(1), "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.UnicastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser malformed_serverID_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", "not_a_number", "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedRequest:  controller.UnicastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser no_UnicastRequest_passed_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.MulticastRequest{},
			expectedRequest:  controller.UnicastRequest{},
			expectedError:    controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser no_namespace_should_fail",
			serverRequest:    &xmlrpc.ServerRequest{"method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			requestToHydrate: &controller.UnicastRequest{},
			expectedError:    controller.FaultDecode.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := UnicastRequestParser(tc.serverRequest, tc.requestToHydrate)
			if err != nil && (tc.expectedError == "" || !strings.Contains(err.Error(), tc.expectedError)) {
				t.Fatalf("expected and actual errors don't match. Expected was:\n%v\nActual is:\n%v:", tc.expectedError, err.Error())
			}
			if err == nil && !reflect.DeepEqual(tc.requestToHydrate, &tc.expectedRequest) {
				t.Fatalf("expected and actual requests don't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedRequest, tc.requestToHydrate)
			}
		})
	}
}
