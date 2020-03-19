package server

import (
	"reflect"
	"testing"
)

func Test_parseToStruct(t *testing.T) {
	type testLoginStruct struct{ Username, Password string }

	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate *testLoginStruct
		expectedStruct  testLoginStruct
		expectedResult  bool
	}{
		{name: "parseToStruct Success",
			args:            []interface{}{"username", "password"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "username", Password: "password"},
			expectedResult:  true},
		{name: "parseToStruct Failed",
			args:            []interface{}{"username", "password"},
			structToHydrate: &testLoginStruct{},
			expectedStruct:  testLoginStruct{Username: "", Password: "password"},
			expectedResult:  false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parseToStruct(tc.args, tc.structToHydrate)
			if reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func Test_parseToMulitcastArgs(t *testing.T) {
	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate *MulticastArgs
		expectedStruct  MulticastArgs
		expectedResult  bool
	}{
		{name: "parseToMulitcastArgs Success",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{HubSessionKey: "sessionKey", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}},
			expectedResult:  true},
		{name: "parseToMulitcastArgs Failed",
			args:            []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHydrate: &MulticastArgs{},
			expectedStruct:  MulticastArgs{HubSessionKey: "", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}},
			expectedResult:  false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parseToMulitcastArgs(tc.args, tc.structToHydrate)
			if reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func Test_parseToUnicastArgs(t *testing.T) {
	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate *UnicastArgs
		expectedStruct  UnicastArgs
		expectedResult  bool
	}{
		{name: "parseToUnicastArgs Success",
			args:            []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{HubSessionKey: "sessionKey", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}},
			expectedResult:  true},
		{name: "parseToUnicastArgs Failed",
			args:            []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHydrate: &UnicastArgs{},
			expectedStruct:  UnicastArgs{HubSessionKey: "", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}},
			expectedResult:  false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parseToUnicastArgs(tc.args, tc.structToHydrate)
			if reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func Test_parseToList(t *testing.T) {
	tt := []struct {
		name            string
		args            []interface{}
		structToHydrate *ListArgs
		expectedStruct  ListArgs
		expectedResult  bool
	}{
		{name: "parseToList Success",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &ListArgs{},
			expectedStruct:  ListArgs{[]interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			expectedResult:  true},
		{name: "parseToList Failed",
			args:            []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHydrate: &ListArgs{},
			expectedStruct:  ListArgs{[]interface{}{"", "arg1_Hub", "arg2_Hub"}},
			expectedResult:  false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parseToList(tc.args, tc.structToHydrate)
			if reflect.DeepEqual(tc.structToHydrate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}
