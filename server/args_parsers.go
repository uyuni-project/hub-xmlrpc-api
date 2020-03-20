package server

import (
	"reflect"
)

type parser func(args []interface{}, output interface{}) error

func parseToStruct(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output)
	for i, arg := range args {
		field := val.Elem().Field(i)
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

func parseToMulitcastArgs(args []interface{}, output interface{}) error {
	hubSessionKey := args[0].(string)
	serverIDs := make([]int64, len(args[1].([]interface{})))
	for i, elem := range args[1].([]interface{}) {
		serverIDs[i] = elem.(int64)
	}

	rest := args[2:len(args)]
	serverArgs := make([][]interface{}, len(rest))

	for i, list := range rest {
		serverArgs[i] = list.([]interface{})
	}

	multicastArgs := []interface{}{hubSessionKey, serverIDs, serverArgs}

	val := reflect.ValueOf(output)
	for i, arg := range multicastArgs {
		field := val.Elem().Field(i)
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

func parseToList(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output)
	if val.Elem().NumField() >= 1 {
		field := val.Elem().Field(0)

		if field.Kind() == reflect.Slice {
			field.Set(reflect.ValueOf(args))
		}
	}
	return nil
}

func parseToUnicastArgs(args []interface{}, output interface{}) error {
	hubSessionKey := args[0].(string)
	serverID := args[1].(int64)

	rest := args[2:len(args)]
	serverArgs := make([]interface{}, len(rest))

	for i, list := range rest {
		serverArgs[i] = list.(interface{})
	}

	unicastArgs := []interface{}{hubSessionKey, serverID, serverArgs}

	val := reflect.ValueOf(output)
	for i, arg := range unicastArgs {
		field := val.Elem().Field(i)
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

//TODO: Temporary exposed for testing
var StructParser = parseToStruct
var UnicastParser = parseToUnicastArgs
var ListParser = parseToList
var MulitcastParser = parseToMulitcastArgs
