# Hub xmlrpc API

Hub XMLRPC API helps to manage servers of servers architecture from a single place. A usual architeture where Hub will be helpful and customer can get benefit from it, looks like this.

![graphviz](https://user-images.githubusercontent.com/12951268/74736042-33518d80-5252-11ea-83a3-04d3d4ae5d11.png)

An Hub user can now manage servers just like as they are normal client of the Uyuni. In addition to that, Hub provides a way to go even further and manage clients
of each server, all from a single place without going to each servers individually.

## Getting Started

These instructions will help user to start with, this section will have some instructiion about the required setup and will give some basic examples which user can modify
according to the needs and play around it.

### Prerequisites

An Hub server should be installed which is in reality a uyuni server with some super powers. At least one uyuni server with one client, should be registered to the Hub to play with the API.



## Installation

Package `hub-xmlrpc-api` is available in Uyuni repo. Refresh the repo if it is not already and install the package. Hub API needs to run on the same machine where Hub is running and will need some additional data. 
It looks for `HUB_CONFIG_FILE` env variable which should be pointing to a json file with the contents like given below.

```json
{
   "type": "json",
    "hub": {
       "manager_api_url": "http://hub-url/rpc/api"
   },
    "connect_timeout": 10,
    "read_write_timeout": 10,
   }
 ```
## Namespaces
Hub supports 3 different namespaces.

1. Hub(default) &#8594; to target hub itself. 
     - `client.hub.login(login, password)`
2. Unicast &#8594; To target a single server
     - `[systems] = client.unicast.system.method_to_call(hubKey, serverId, arg1, arg2)`
3. Multicast &#8594; To target multiple servers
     - `[ [systems], [systems] ] = client.multicast.system.list_systems(hubKey, [serverIds],[arg1..], [arg2..])`

Please note the return type here, in case of unicast, data comes back in an array while in case of multicast as array of array.     
##Note : As hub is a uyuni server itself, user can call the existing uyuni xmlrpc api using the same connection instead of a seprate connection and Hub API will take care of rest.
## Modes
Hub supports 3 different modes

1. manual mode(default) &#8594; User needs to provide  **credentials ** for each server she wants to perform operation
2. relay mode  &#8594; Use Hub  **credentials ** but  **Servers** list will be provided by the user
3. auto connect mode &#8594; Use Hub  **credentials ** +  **Servers** list will detected  **automatically** where Hub user has
access to

An example is given below to help it more clear. Here the auto connect mode is enabled and other modes are commented out. 

```python

#!/usr/bin/python
import xmlrpclib
import itertools
 
HUB_URL = "http://localhost:8000/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"
 
client = xmlrpclib.Server(HUB_URL, verbose=0)

#login
#hubKey = client.hub.login(HUB_LOGIN, HUB_PASSWORD )
#hubKey = client.hub.LoginWithAuthRelayMode(HUB_LOGIN, HUB_PASSWORD )
#hubKey = client.hub.LoginWithAutoconnectMode(HUB_LOGIN, HUB_PASSWORD )

#get list of server ids registerd with hub
serverIds = client.hub.listServerIds(hubKey)


#Manual-------------
#usernames = ["admin" for s in serverIds]
#passwords = ["admin" for s in serverIds]
#client.hub.attachToServers(hubKey, serverIds, usernames, passwords)


#Relay_mode---------
#client.hub.attachToServers(hubKey, serverIds)


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

* **####** - *Initial work* 

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the --- License - see the [LICENSE.md](LICENSE.md) file for details


