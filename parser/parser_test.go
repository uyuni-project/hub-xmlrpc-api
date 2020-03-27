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
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  testLoginStruct
		expectedMethod  string
		expectedError   string
	}{
		{name: "parseToStruct Success",
			args:            []interface{}{"username", "password"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "username", Password: "password"},
			expectedMethod:  "method"},
		{name: "parseToStruct no_struct_passed Failed",
			args:            []interface{}{"username", "password"},
			structToHydrate: &[]interface{}{},
			expectedStruct:  testLoginStruct{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToStruct wrong_number_of_arguments_passed Failed",
			args:            []interface{}{"username", "password", "extra_argument"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedMethod:  "method",
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToStruct wrong_type_of_arguments_passed Failed",
			args:            []interface{}{"username", 123},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToStruct(tc.expectedMethod, tc.args, tc.structToHydrate)
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
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  server.ListArgs
		expectedMethod  string
		expectedError   string
	}{
		{name: "parseToList Success",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &server.ListArgs{},
			expectedStruct:  server.ListArgs{Method: "method", Args: []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			expectedMethod:  "method"},
		{name: "parseToList no_ListArgs_passed Failed",
			args:            []interface{}{},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.ListArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToList(tc.expectedMethod, tc.args, tc.structToHydrate)
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
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  server.MulticastArgs
		expectedMethod  string
		expectedError   string
	}{
		{name: "parseToMulitcastArgs Success",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{Method: "method", HubSessionKey: "sessionKey", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}},
			expectedMethod:  "method"},
		{name: "parseToMulitcastArgs no_serverID_passed Failed",
			args:            []interface{}{"sessionKey"},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToMulitcastArgs malformed_hubSessionKey Failed",
			args:            []interface{}{123, []interface{}{int64(1000010001)}, []interface{}{"arg1_Server1"}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToMulitcastArgs malformed_serverID Failed",
			args:            []interface{}{"sessionKey", []interface{}{"1000010001", "1000010002"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToMulitcastArgs no_MulticastArgs_passed Failed",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001)}, "arg1_Server1", "arg2_Server1"},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.MulticastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToMulitcastArgs(tc.expectedMethod, tc.args, tc.structToHydrate)
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
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  server.UnicastArgs
		expectedMethod  string
		expectedError   string
	}{
		{name: "parseToUnicastArgs Success",
			args:            []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{Method: "method", HubSessionKey: "sessionKey", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}},
			expectedMethod:  "method"},
		{name: "parseToUnicastArgs wrong_number_of_arguments Failed",
			args:            []interface{}{"sessionKey"},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultWrongArgumentsNumber.Message},
		{name: "parseToUnicastArgs malformed_hubSessionKey Failed",
			args:            []interface{}{123, "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
		{name: "parseToUnicastArgs malformed_serverId Failed",
			args:            []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &server.UnicastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedError:   codec.FaultInvalidParams.Message,
			expectedMethod:  "method"},
		{name: "parseToUnicastArgs no_UnicastArgs_passed Failed",
			args:            []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &server.MulticastArgs{},
			expectedStruct:  server.UnicastArgs{},
			expectedMethod:  "method",
			expectedError:   codec.FaultInvalidParams.Message},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToUnicastArgs(tc.expectedMethod, tc.args, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}
