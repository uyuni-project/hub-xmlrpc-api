# Hub XMLRPC API

The Hub XMLRPC API helps to operate SUSE Manager or Uyuni infrastructures with one Server, called a Hub, managing several Servers. The typical architecture is represented below:

![graphviz](https://user-images.githubusercontent.com/12951268/74736042-33518d80-5252-11ea-83a3-04d3d4ae5d11.png)

Servers are registered and can be managed from the Hub just like normal SUSE Manager or Uyuni clients. In addition to that, the Hub XMLRPC API provides a more convenient way to access all Servers' XMLRPC APIs from one only endpoint in the Hub.

## Getting Started

### Requirements
 - one central SUSE Manager/Uyuni Server (hereinafter: "Hub")
 - two or more peripheral Servers registered to the Hub as Salt clients
 - any number of end clients registered to peripheral Servers

You will need credentials to all Server XMLRPC APIs, including the Hub's.

### Installation

Install the package `hub-xmlrpc-api`, available from SUSE Manager 4.1 and Uyuni repos, either on the Hub itself or on a host that has access to all Servers' XMLRPC APIs.

Configuration of `hub-xmlrpc-api` is specified in a JSON file like the following:

```json
{
   "type": "json",
    "hub": {
       "manager_api_url": "http://localhost/rpc/api"
   },
    "connect_timeout": 10,
    "read_write_timeout": 10,
   }
 ```

Replace `localhost` in `http://localhost/rpc/api` above with the Hub's FQDN if necessary.

Set the `HUB_CONFIG_FILE` environment variable to point to the configuration file.

## Usage

`hub-xmlrpc-api` is a daemon, currently to be launched from the command line.


Once running, you can connect to the `hub-xmlrpc-api` at port 8888 via any XMLRPC compliant client libraries (see examples below).


### Namespaces

Hub supports 3 different namespaces.

1. `hub` &#8594; to target the Hub itself:
     - example: `hubKey = client.hub.login(login, password)`
2. `unicast` &#8594; to target a single Server registered in the Hub:
     - example: `systems = client.unicast.system.list_systems(hubKey, serverId)`
3. `multicast` &#8594; to target multiple Servers
     - example: `systems_per_server = client.multicast.system.list_systems(hubKey, [serverId1, serverId2], [], [])`

Note that:
 - all XMLRPC API methods available in a single Server are exposed by the namespaces above. Generally speaking, they accept the same parameters and return the same values with the exceptions described below
 - the `hubKey` can be obtained via the `client.hub.login(login, password)` method
 - individual Server IDs can be obtained via `client.hub.listServerIds(hubKey)` (see example below)
 - the `unicast` namespace assumes all methods receive `hubKey` and `serverId` as their first two parameters, then any other parameter as specified by the regular Server API
 - the `multicast` namespace assumes all methods receive `hubKey`, a list of Server IDs, then lists of per-Server parameters as specified by the regular Server API. Return value will be an array, indexed per Server, of the results of individual Server calls

### Authentication modes

Hub supports 3 different authentication modes.

1. manual mode (default): user needs to provide API credentials for each Server explicitly
2. relay mode: the same credentials used to authenticate against the Hub will be re-used to authenticate Servers. The list of Servers to connect to will still be provided by the user
3. auto connect mode: Hub credentials will be reused for Servers and any Server the user has access to will be automatically connected

### Python example

```python

#!/usr/bin/python
import xmlrpclib
import itertools
 
HUB_URL = "http://localhost:8888/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"
 
client = xmlrpclib.Server(HUB_URL, verbose=0)

# Login (uncomment one line)
#hubKey = client.hub.login(HUB_LOGIN, HUB_PASSWORD)
#hubKey = client.hub.LoginWithAuthRelayMode(HUB_LOGIN, HUB_PASSWORD)
#hubKey = client.hub.LoginWithAutoconnectMode(HUB_LOGIN, HUB_PASSWORD)

# get list of Server ids registerd to the Hub
serverIds = client.hub.listServerIds(hubKey)

# Manual authentication mode example. Uncomment if `login` was uncommented above
#usernames = ["admin" for s in serverIds]
#passwords = ["admin" for s in serverIds]
#client.hub.attachToServers(hubKey, serverIds, usernames, passwords)

# Relay authentication mode example. Uncomment if `LoginWithAuthRelayMode` was uncommented above
#client.hub.attachToServers(hubKey, serverIds)

# Nothing has to be done if `LoginWithAutoconnectMode` was uncommented above

# Run call
systemsPerServer = client.multicast.system.list_systems(hubKey, serverIds)

for system in itertools.chain.from_iterable(systemsPerServer):
  print system.get('name')

#logout
client.auth.logout(hubKey)
```


## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

See the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.


