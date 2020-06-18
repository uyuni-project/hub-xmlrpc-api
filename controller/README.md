Controllers are called by the Codec Gorilla component and marshal/unmarshal requests and responses to functions in the gaeway package. There is one Controller per supported method or namespace.

Functions in the gateway package actually implement that method or namespace's functionality.