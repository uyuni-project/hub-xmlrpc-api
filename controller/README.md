Controllers are called by the Codec Gorilla component and marshal/unmarshal requests and responses to functions in the gateway package. There is one Controller per supported usecase.

They provide an entry point for the application, and expose the available usecases to the user.
They should not contain any bussiness logic at all, rather they should delegate this to the "gateway" package, which contains the implementation of the usecases.

The responsibilities of the Controllers are:
- expose the available usecases to the user.
- validate the user request data.
- delegate the usecase execution to the "gateway" package.
- transform the user request to the proper data structure require by the usecases, implemented in the "gateway" package.
- transform the response provided by the execution of the usecases to the proper model required by the user. 