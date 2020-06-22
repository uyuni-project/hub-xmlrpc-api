Controllers are called by the Codec Gorilla component and marshal/unmarshal requests and responses to functions in the gateway package. There is one Controller per supported method or namespace.

Functions in the gateway package actually implement that method or namespace's functionality, in other words, the "business rules" or "use cases" minus protocol and other low-level details.