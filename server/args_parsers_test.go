package server

import (
	"reflect"
	"testing"
)

func TestStructParser(t *testing.T) {
	type testLoginStruct struct{ Username, Password string }

	tt := []struct {
		name              string
		args              []interface{}
		structToHidratate *testLoginStruct
		expectedStruct    testLoginStruct
		expectedResult    bool
	}{
		{name: "StructParser Success",
			args:              []interface{}{"username", "password"},
			structToHidratate: &testLoginStruct{},
			expectedStruct:    testLoginStruct{Username: "username", Password: "password"},
			expectedResult:    true},
		{name: "StructParser Failed",
			args:              []interface{}{"username", "password"},
			structToHidratate: &testLoginStruct{},
			expectedStruct:    testLoginStruct{Username: "", Password: "password"},
			expectedResult:    false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			new(StructParser).Parse(tc.args, tc.structToHidratate)
			if reflect.DeepEqual(tc.structToHidratate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func TestMulticastArgsParser(t *testing.T) {
	tt := []struct {
		name              string
		args              []interface{}
		structToHidratate *MulticastArgs
		expectedStruct    MulticastArgs
		expectedResult    bool
	}{
		{name: "MuticastArgsParser Success",
			args:              []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHidratate: &MulticastArgs{},
			expectedStruct:    MulticastArgs{HubSessionKey: "sessionKey", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}},
			expectedResult:    true},
		{name: "MuticastArgsParser Failed",
			args:              []interface{}{"sessionKey", []interface{}{int64(1000010001), int64(1000010002)}, []interface{}{"arg1_Server1", "arg1_Server2"}, []interface{}{"arg2_Server1", "arg2_Server2"}},
			structToHidratate: &MulticastArgs{},
			expectedStruct:    MulticastArgs{HubSessionKey: "", ServerIDs: []int64{1000010001, 1000010002}, ServerArgs: [][]interface{}{{"arg1_Server1", "arg1_Server2"}, {"arg2_Server1", "arg2_Server2"}}},
			expectedResult:    false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			new(MulticastArgsParser).Parse(tc.args, tc.structToHidratate)
			if reflect.DeepEqual(tc.structToHidratate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func TestUnicastArgsParser(t *testing.T) {
	tt := []struct {
		name              string
		args              []interface{}
		structToHidratate *UnicastArgs
		expectedStruct    UnicastArgs
		expectedResult    bool
	}{
		{name: "UnicastArgsParser Success",
			args:              []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHidratate: &UnicastArgs{},
			expectedStruct:    UnicastArgs{HubSessionKey: "sessionKey", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}},
			expectedResult:    true},
		{name: "UnicastArgsParser Failed",
			args:              []interface{}{"sessionKey", int64(1000010001), "arg1_Server1", "arg2_Server1"},
			structToHidratate: &UnicastArgs{},
			expectedStruct:    UnicastArgs{HubSessionKey: "", ServerID: int64(1000010001), ServerArgs: []interface{}{"arg1_Server1", "arg2_Server1"}},
			expectedResult:    false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			new(UnicastArgsParser).Parse(tc.args, tc.structToHidratate)
			if reflect.DeepEqual(tc.structToHidratate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}

func TestListParser(t *testing.T) {
	tt := []struct {
		name              string
		args              []interface{}
		structToHidratate *ListArgs
		expectedStruct    ListArgs
		expectedResult    bool
	}{
		{name: "ListParser Success",
			args:              []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHidratate: &ListArgs{},
			expectedStruct:    ListArgs{[]interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"}},
			expectedResult:    true},
		{name: "ListParser Failed",
			args:              []interface{}{"sessionKey", "arg1_Hub", "arg2_Hub"},
			structToHidratate: &ListArgs{},
			expectedStruct:    ListArgs{[]interface{}{"", "arg1_Hub", "arg2_Hub"}},
			expectedResult:    false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			new(ListParser).Parse(tc.args, tc.structToHidratate)
			if reflect.DeepEqual(tc.structToHidratate, &tc.expectedStruct) != tc.expectedResult {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
			}
		})
	}
}
