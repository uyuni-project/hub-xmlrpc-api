This package contains the bussines rules of the application, and the implementation of all the usecases.
All the other components of the application are plugins to this package.
It should not contain any external dependency at all, and it should interact with the external packages only through interfaces.

At the moment, the following use cases are implemented:

- "Login": login to the Hub server, using credentials provided by the user. 
- "Login with auth-relay mode": Extends the "Login" use case, and stores the credentials for further usage to login to the Uyuni peripheral servers.
- "Login with Autoconnect mode": Extends the "Login with auth-realy mode" and includes the "Attach to servers" usecase.
- "Attach to servers": login to the Uyuni peripheral servers of a have architecture.
- "Logout": Logout from the Hub server and the Uyuni peripheral servers.
- "Unicast": Execute a unicast call to a peripheral Uyuni server of the Hub architecture.
- "Multicast": Execute a multicast call to many peripheral Uyuni servers of the Hub architecture.
- "Proxy call to Hub": Redirect a call to the underlying regular Uyuni XMLPRC API of the Hub server.
- "List user server IDs": Retrieve a the list of the server IDs of those Uyuni peripheral servers a user has access to.