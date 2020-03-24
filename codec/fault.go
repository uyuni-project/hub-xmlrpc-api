package codec

import (
	"fmt"
)

var (
	FaultInvalidParams        = Fault{Code: -32602, Message: "Invalid Method Parameters"}
	FaultWrongArgumentsNumber = Fault{Code: -32602, Message: "Wrong Arguments Number"}
	FaultInternalError        = Fault{Code: -32603, Message: "Internal Server Error"}
	FaultApplicationError     = Fault{Code: -32500, Message: "Application Error"}
	FaultSystemError          = Fault{Code: -32400, Message: "System Error"}
	FaultDecode               = Fault{Code: -32700, Message: "Parsing error: not well formed"}
	FaultInvalidCredentials   = Fault{Code: 2950, Message: "Either the password or username is incorrect"}
)

type Fault struct {
	Code    int    `xmlrpc:"faultCode"`
	Message string `xmlrpc:"faultString"`
}

func (f Fault) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.Message)
}
