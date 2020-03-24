package server

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseToStruct(t *testing.T) {
	type testLoginStruct struct{ Username, Password string }

	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  testLoginStruct
		expectedError   string
	}{
		{name: "parseToStruct Success",
			args:            []interface{}{"username", "password"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "username", Password: "password"}},
		{name: "parseToStruct no_struct_passed Failed",
			args:            []interface{}{"username", "password"},
			structToHydrate: &[]interface{}{},
			expectedStruct:  testLoginStruct{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToStruct wrong_number_of_arguments_passed Failed",
			args:            []interface{}{"username", "password", "extra_argument"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   FaultWrongArgumentsNumber.String},
		{name: "parseToStruct wrong_type_of_arguments_passed Failed",
			args:            []interface{}{"username", 123},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{},
			expectedError:   FaultInvalidParams.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToStruct(tc.args, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}

func Test_parseToList(t *testing.T) {
	var emptyStruct struct{}

	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate interface{}
		expectedStruct  ListArgs
		expectedError   string
	}{
		{name: "parseToList Success",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &ListArgs{},
			expectedStruct:  ListArgs{[]interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}}},
		{name: "parseToList no_struct_passed Failed",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &[]interface{}{},
			expectedStruct:  ListArgs{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToList wrong_number_of_arguments_passed Failed",
			args:            []interface{}{"username", "password", "extra_argument"},
			structToHydrate: &emptyStruct,
			expectedStruct:  ListArgs{},
			expectedError:   FaultWrongArgumentsNumber.String},
		{name: "parseToList no_list_field_in_struct Failed",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  ListArgs{},
			expectedError:   FaultInvalidParams.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToList(tc.args, tc.structToHydrate)
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
		expectedStruct  MulticastArgs
		expectedError   string
	}{
		{name: "parseToMulitcastArgs Success",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{HubSessionKey: "sessionKey", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}}},
		{name: "parseToMulitcastArgs no_serverID_passed Failed",
			args:            []interface{}{"sessionKey"},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{},
			expectedError:   FaultWrongArgumentsNumber.String},
		{name: "parseToMulitcastArgs malformed_hubSessionKey Failed",
			args:            []interface{}{123, []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToMulitcastArgs malformed_serverID Failed",
			args:            []interface{}{"sessionKey", []interface{}{"1000010001", "1000010002"}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToMulitcastArgs no_MulticastArgs_passed Failed",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001)}, "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  MulticastArgs{},
			expectedError:   FaultInvalidParams.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToMulitcastArgs(tc.args, tc.structToHydrate)
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
		expectedStruct  UnicastArgs
		expectedError   string
	}{
		{name: "parseToUnicastArgs Success",
			args:            []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{HubSessionKey: "sessionKey", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}}},
		{name: "parseToUnicastArgs wrong_number_of_arguments Failed",
			args:            []interface{}{"sessionKey"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{},
			expectedError:   FaultWrongArgumentsNumber.String},
		{name: "parseToUnicastArgs malformed_hubSessionKey Failed",
			args:            []interface{}{123, "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToUnicastArgs malformed_serverId Failed",
			args:            []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{},
			expectedError:   FaultInvalidParams.String},
		{name: "parseToUnicastArgs no_UnicastArgs_passed Failed",
			args:            []interface{}{"sessionKey", "1000010001", "arg1_Server1", "arg2_Server1"},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  UnicastArgs{},
			expectedError:   FaultInvalidParams.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := parseToUnicastArgs(tc.args, tc.structToHydrate)
			if err != nil && !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}

			if err == nil && !reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedStruct)
			}
		})
	}
}
