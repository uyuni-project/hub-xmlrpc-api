package codec

import (
	"fmt"
)

var (
	FaultInvalidParams        = FaultError{Code: -32602, Message: "Invalid Method Parameters"}
	FaultWrongArgumentsNumber = FaultError{Code: -32602, Message: "Wrong Arguments Number"}
	FaultInternalError        = FaultError{Code: -32603, Message: "Internal Server Error"}
	FaultApplicationError     = FaultError{Code: -32500, Message: "Application Error"}
	FaultSystemError          = FaultError{Code: -32400, Message: "System Error"}
	FaultDecode               = FaultError{Code: -32700, Message: "Parsing error: not well formed"}
	FaultInvalidCredentials   = FaultError{Code: 2950, Message: "Either the password or username is incorrect"}
)

type FaultError struct {
	Code    int    `xmlrpc:"faultCode"`
	Message string `xmlrpc:"faultString"`
}

func (f FaultError) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.Message)
}
