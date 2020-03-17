package server

import (
	"reflect"
)

type Parser interface {
	Parse(args []interface{}, output interface{}) error
}

type StructParser struct{}

func (p *StructParser) Parse(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output)
	for i, arg := range args {
		field := val.Elem().Field(i)
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

type MulticastArgsParser struct{}

func (p *MulticastArgsParser) Parse(args []interface{}, output interface{}) error {
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

type ListParser struct{}

func (p *ListParser) Parse(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output)
	if val.Elem().NumField() >= 1 {
		field := val.Elem().Field(0)

		if field.Kind() == reflect.Slice {
			field.Set(reflect.ValueOf(args))
		}
	}
	return nil
}

type UnicastArgsParser struct{}

func (p *UnicastArgsParser) Parse(args []interface{}, output interface{}) error {
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
