package server

import (
	"log"
	"reflect"
)

var (
	StructParser    = parseToStruct
	UnicastParser   = parseToUnicastArgs
	ListParser      = parseToList
	MulitcastParser = parseToMulitcastArgs
)

type parser func(args []interface{}, output interface{}) error

func parseToStruct(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output).Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("Error ocurred when parsing arguments")
		return FaultInvalidParams
	}
	if val.NumField() < len(args) {
		log.Printf("Error ocurred when parsing arguments")
		return FaultWrongArgumentsNumber
	}

	for i, arg := range args {
		field := val.Field(i)
		if field.Type() != reflect.ValueOf(arg).Type() {
			log.Printf("Error ocurred when parsing arguments")
			return FaultInvalidParams
		}
		field.Set(reflect.ValueOf(arg))
	}
	return nil
}

func parseToUnicastArgs(args []interface{}, output interface{}) error {
	parsedArgs, ok := output.(*UnicastArgs)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return FaultInvalidParams
	}
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return FaultInvalidParams
	}

	serverID, ok := args[1].(int64)
	if !ok {
		log.Printf("Error ocurred when parsing serverID argument")
		return FaultInvalidParams
	}

	rest := args[2:len(args)]
	serverArgs := make([]interface{}, len(rest))
	for i, list := range rest {
		serverArgs[i] = list.(interface{})
	}

	*parsedArgs = UnicastArgs{hubSessionKey, serverID, serverArgs}
	return nil
}

func parseToMulitcastArgs(args []interface{}, output interface{}) error {
	parsedArgs, ok := output.(*MulticastArgs)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return FaultInvalidParams
	}
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return FaultInvalidParams
	}

	serverIDs := make([]int64, len(args[1].([]interface{})))
	for i, elem := range args[1].([]interface{}) {
		serverIDs[i], ok = elem.(int64)
		if !ok {
			log.Printf("Error ocurred when parsing serverIDs argument")
			return FaultInvalidParams
		}
	}

	rest := args[2:len(args)]
	serverArgs := make([][]interface{}, len(rest))
	for i, list := range rest {
		serverArgs[i] = list.([]interface{})
	}

	*parsedArgs = MulticastArgs{hubSessionKey, serverIDs, serverArgs}
	return nil
}

func parseToList(args []interface{}, output interface{}) error {
	val := reflect.ValueOf(output).Elem()
	if val.Kind() != reflect.Struct {
		log.Printf("Error ocurred when parsing arguments")
		return FaultInvalidParams
	}
	if val.NumField() < 1 {
		log.Printf("Error ocurred when parsing arguments")
		return FaultWrongArgumentsNumber
	}

	field := val.Field(0)
	if field.Kind() != reflect.Slice {
		log.Printf("Error ocurred when parsing arguments")
		return FaultInvalidParams
	}

	field.Set(reflect.ValueOf(args))
	return nil
}
