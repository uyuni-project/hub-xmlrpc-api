package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/codec"
)

func Test_parseToStruct(t *testing.T) {
	type testLoginStruct struct{ Username, Password string }

	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  testLoginStruct
		expectedError   string
	}{
		{name: "parseToStruct Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", "password"}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "username", Password: "password"}},
		{name: "parseToStruct no_struct_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", "password"}},
			structToHydrate: &[]interface{}{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToStruct wrong_number_of_arguments_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", "password", "extra_argument"}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultWrongArgumentsNumber.Message},
		{name: "parseToStruct wrong_type_of_arguments_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", 123}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToStruct(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToListRequest(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.ListRequest
		expectedError   string
	}{
		{name: "parseToListRequest Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			structToHydrate: &controller.ListRequest{},
			expectedStruct:  controller.ListRequest{Method: "method", Args: []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}}},
		{name: "parseToListRequest no_ListRequest_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.ListRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToListRequest(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToMulitcastRequest(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.MulticastRequest
		expectedError   string
	}{
		{name: "parseToMulitcastRequest Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{int64(1), int64(2)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{Method: "method", HubSessionKey: "hubSessionKey", ArgsByServer: map[int64][]interface{}{1: []interface{}{"arg1_Server1", "arg2_Server1"}, 2: []interface{}{"arg1_Server2", "arg2_Server2"}}}},
		{name: "parseToMulitcastRequest no_serverID_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{"serverSessionKey"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToMulitcastRequest malformed_hubSessionKey Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{1, []interface{}{int64(2)}, []interface{}{"arg1_Server1"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToMulitcastRequest malformed_serverID Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{"1", "2"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToMulitcastRequest no_MulticastRequest_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"hubSessionKey", []interface{}{int64(1)}, "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.MulticastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToMulitcastRequest(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToUnicastRequest(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  controller.UnicastRequest
		expectedError   string
	}{
		{name: "parseToUnicastRequest Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", int64(1), "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{Method: "method", HubSessionKey: "sessionKey", ServerID: int64(1), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}}},
		{name: "parseToUnicastRequest wrong_number_of_arguments Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultWrongArgumentsNumber.Message},
		{name: "parseToUnicastRequest malformed_hubSessionKey Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{123, "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToUnicastRequest malformed_serverId Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.UnicastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
		{name: "parseToUnicastRequest no_UnicastRequest_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "1", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &controller.MulticastRequest{},
			expectedStruct:  controller.UnicastRequest{},
			expectedError:   controller.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToUnicastRequest(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}
