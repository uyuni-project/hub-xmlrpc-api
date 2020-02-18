# Hub xmlrpc API

A Hub is an Uyuni Server that can manage other Uyuni Servers registered to it. It can be seen as a "server of servers".

Hub XMLRPC API helps to manage this architecture from a single place. A usual architeture where Hub will be helpful and customer can get benefit from it, looks like this:

![graphviz](https://user-images.githubusercontent.com/12951268/74736042-33518d80-5252-11ea-83a3-04d3d4ae5d11.png)

Users defined on the Hub can now manage Uyuni Servers registered to it as regular systems through the Hub XMLRPC API. In addition to that, the Hub XMLRPC API provides a way to go even further and manage the systems registered to each SUMA Server, all from a single place without the need to directly accessing to each server individually.

## Getting Started

This section is a quick start guide with instructions to setup the Hub XMLRPC API. It also contains some basic examples scripts to work with, and instructions to setup a Hub architechture, as describe in https://github.com/uyuni-project/hub-xmlrpc-api/tree/master#hub-xmlrpc-api .

### Setting up the environment

To work with the Hub XMLRPC API you need a working Hub environment, as described in https://github.com/uyuni-project/hub-xmlrpc-api/tree/master#hub-xmlrpc-api . You can set up such installation by manual installation and registration of Uyuni Servers, or you can use Sumaform (https://github.com/uyuni-project/sumaform) for it, by adding the following modules in your main.tf file:

```yaml
module "hub-server" {
  source = "./modules/libvirt/suse_manager"
  base_configuration = "${module.base.configuration}"
  product_version = "uyuni-master"
  name = "hub-server"
  image = "opensuse151"
  use_os_released_updates = true
}

module "slave-server-1" {
  source = "./modules/libvirt/suse_manager"
  base_configuration = "${module.base.configuration}"
  product_version = "uyuni-master"
  name = "slave-server-1"
  image = "opensuse151"
  register_to_server = "mch-hub-server.tf.local"
}

module "min-for-slave-server-1" {
  source = "./modules/libvirt/minion"
  base_configuration = "${module.base.configuration}"
  product_version = "uyuni-master"
  name = "min-for-slave-server-1"
  image = "sles12sp3"
  server_configuration = "${module.slave-server-1.configuration}"
}

module "slave-server-2" {
  source = "./modules/libvirt/suse_manager"
  base_configuration = "${module.base.configuration}"
  product_version = "uyuni-master"
  name = "slave-server-2"
  image = "opensuse151"
  register_to_server = "mch-hub-server.tf.local"
}

module "min-for-slave-server-2" {
  source = "./modules/libvirt/minion"
  base_configuration = "${module.base.configuration}"
  product_version = "uyuni-master"
  name = "min-for-slave-server-1"
  image = "sles12sp3"
  server_configuration = "${module.slave-server-2.configuration}"
}
```

This will create a Hub Uyuni Server (hub-server) with two slave Uyuni Servers (slave-server-1 and slave-server-2) registered to the Hub, and a minion registered to each slave Uyuni Server (min-for-slave-server-1 and min-for-slave-server-2).

### Prerequisites

An Hub server should be installed which is in reality a uyuni server with some super powers. At least one uyuni server with one client, should be registered to the Hub to play with the API.



## Installation of the Hub XMLRPC API

Package `hub-xmlrpc-api` is available in the Uyuni repository. Refresh the repository if it is not already and install the package. The Hub XMLRPC API needs to run on the same machine where Hub Uyuni Server is running and will need some additional data. 
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
Hub supports 3 different namespaces:

1. Hub(default) &#8594; to target the hub itself:
     - `client.hub.login(login, password)`
2. Unicast &#8594; To target a single slave server:
     - `[systems] = client.unicast.system.method_to_call(hubKey, serverId, arg1, arg2)`
3. Multicast &#8594; To target multiple slave servers:
     - `{ "Successfull": {"Responses": [systems], "ServerIds":[]}, "Failed": {"Responses": [systems], "ServerIds":[]} } = client.multicast.system.list_systems(hubKey, [serverIds],[arg1..], [arg2..])`

Please note the return type here, in case of unicast, data comes back in an array while in case of multicast the response format is a JSON object, with to fields: "Successfull" and "Failed". Each of these fields contain another JSON object with the fields "Responses" (the servers responses) and "ServerIds" (the Ids of the servers whose response is in the "Responses" field). "Responses" and "ServerIds" come in matchin order, meaning the first Response matches the first ServerId, and so on.
##Note : As hub is an Uyuni Server itself, the user can call the existing Uyuni XMLRPC API using the same connection instead of a separate connection, and the Hub XMLRP API will take care of rest.
## Modes
Hub supports 3 different modes:

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


