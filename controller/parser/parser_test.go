package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func Test_LoginRequestParser(t *testing.T) {
	type testLoginStruct struct{ Username, Password string }

	tt := []struct {
		name            string
		serverRequest   *xmlrpc.ServerRequest
		structToHydrate interface{}
		expectedStruct  testLoginStruct
		expectedError   string
	}{
		{name: "LoginRequestParser Success",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"username", "password"}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "username", Password: "password"}},
		{name: "LoginRequestParser no_struct_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"username", "password"}},
			structToHydrate: &[]interface{}{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "LoginRequestParser wrong_number_of_arguments_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"username", "password", "extra_argument"}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultWrongArgumentsNumber.Message},
		{name: "LoginRequestParser wrong_type_of_arguments_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"username", 123}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := LoginRequestParser(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedStruct, tc.structToHydrate)
			}
		})
	}
}

func Test_ProxyCallToHubRequestParser(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *xmlrpc.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.ProxyCallToHubRequest
		expectedError   string
	}{
		{name: "ProxyCallToHubRequestParser Success",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			structToHydrate: &controller.ProxyCallToHubRequest{},
			expectedStruct:  controller.ProxyCallToHubRequest{Call: "method", Args: []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}}},
		{name: "ProxyCallToHubRequestParser no_ProxyCallToHubRequest_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.ProxyCallToHubRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := ProxyCallToHubRequestParser(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedStruct, tc.structToHydrate)
			}
		})
	}
}

func Test_MulticastRequestParser(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *xmlrpc.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.MulticastRequest
		expectedError   string
	}{
		{name: "MulticastRequestParser Success",
			serverRequest:   &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{Call: "method", HubSessionKey: "hubSessionKey", ServerIDs: []int64{1, 2}, ArgsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}}}},
		{name: "MulticastRequestParser no_serverID_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{"serverSessionKey"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser malformed_hubSessionKey Failed",
			serverRequest:   &xmlrpc.ServerRequest{"multicast.method", []interface{}{1, []interface{}{int64(2)}, []interface{}{"arg1_Server1"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser malformed_serverID Failed",
			serverRequest:   &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{"1", "2"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser no_MulticastRequest_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"multicast.method", []interface{}{"hubSessionKey", []interface{}{int64(1)}, "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "MulticastRequestParser no_namespace Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedError:   controller.FaultDecode.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := MulticastRequestParser(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedStruct, tc.structToHydrate)
			}
		})
	}
}

func Test_UnicastRequestParser(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *xmlrpc.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.UnicastRequest
		expectedError   string
	}{
		{name: "UnicastRequestParser Success",
			serverRequest:   &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{Call: "method", HubSessionKey: "sessionKey", ServerID: int64(1), Args: []interface{}{"arg1_Server1", "arg2_Server1"}}},
		{name: "UnicastRequestParser wrong_number_of_arguments Failed",
			serverRequest:   &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultWrongArgumentsNumber.Message},
		{name: "UnicastRequestParser malformed_hubSessionKey Failed",
			serverRequest:   &xmlrpc.ServerRequest{"unicast.method", []interface{}{123, "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser malformed_serverId Failed",
			serverRequest:   &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser no_UnicastRequest_passed Failed",
			serverRequest:   &xmlrpc.ServerRequest{"unicast.method", []interface{}{"sessionKey", "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "UnicastRequestParser no_namespace Failed",
			serverRequest:   &xmlrpc.ServerRequest{"method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedError:   controller.FaultDecode.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := UnicastRequestParser(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", &tc.expectedStruct, tc.structToHydrate)
			}
		})
	}
}
