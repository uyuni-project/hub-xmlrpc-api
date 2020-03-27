package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/codec"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
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
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToStruct wrong_number_of_arguments_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", "password", "extra_argument"}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToStruct wrong_type_of_arguments_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"username", 123}},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   codec.FaultInvalidParams.Message},
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

func Test_parseToList(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  server.ListArgs
		expectedError   string
	}{
		{name: "parseToList Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			structToHydrate: &server.ListArgs{},
			expectedStruct:  server.ListArgs{Method: "method", Args: []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}}},
		{name: "parseToList no_ListArgs_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.ListArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToList(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToMulitcastArgs(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  server.MulticastArgs
		expectedError   string
	}{
		{name: "parseToMulitcastArgs Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{Method: "method", HubSessionKey: "sessionKey", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}}},
		{name: "parseToMulitcastArgs no_serverID_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", []interface{}{"sessionKey"}}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToMulitcastArgs malformed_hubSessionKey Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{123, []interface{}{int64(1000010001)}, []interface{}{"arg1_Server1"}}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToMulitcastArgs malformed_serverID Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", []interface{}{"1000010001", "1000010002"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToMulitcastArgs no_MulticastArgs_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", []interface{}{int64(1000010001)}, "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToMulitcastArgs(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToUnicastArgs(t *testing.T) {
	tt := []struct {
		name            string
		serverRequest   *codec.ServerRequest
		structToHydrate interface{}
		expectedStruct  server.UnicastArgs
		expectedError   string
	}{
		{name: "parseToUnicastArgs Success",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{Method: "method", HubSessionKey: "sessionKey", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}}},
		{name: "parseToUnicastArgs wrong_number_of_arguments Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey"}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToUnicastArgs malformed_hubSessionKey Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{123, "1000010001", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToUnicastArgs malformed_serverId Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToUnicastArgs no_UnicastArgs_passed Failed",
			serverRequest:   &codec.ServerRequest{"method", []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToUnicastArgs(tc.serverRequest, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}
