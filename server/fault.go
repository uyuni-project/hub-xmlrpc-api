package server

import (
	"fmt"

	"github.com/kolo/xmlrpc"
)

var (
	FaultInvalidParams        = Fault{Code: -32602, String: "Invalid Method Parameters"}
	FaultWrongArgumentsNumber = Fault{Code: -32602, String: "Wrong Arguments Number"}
	FaultInternalError        = Fault{Code: -32603, String: "Internal Server Error"}
	FaultApplicationError     = Fault{Code: -32500, String: "Application Error"}
	FaultSystemError          = Fault{Code: -32400, String: "System Error"}
	FaultDecode               = Fault{Code: -32700, String: "Parsing error: not well formed"}
	FaultInvalidCredentials   = Fault{Code: 2950, String: "Either the password or username is incorrect."}
)

type Fault struct {
	Code   int    `xmlrpc:"faultCode"`
	String string `xmlrpc:"faultString"`
}

func (f Fault) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

func fault2XML(fault Fault) string {
	buffer := "<methodResponse><fault>"
	xmlByte, _ := xmlrpc.Marshal(fault)
	buffer += string(xmlByte)
	buffer += "</fault></methodResponse>"
	return buffer
}
