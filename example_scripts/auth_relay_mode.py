#!/usr/bin/python
import xmlrpclib  
import itertools

HUB_URL = "http://localhost:8888/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"

client = xmlrpclib.Server(HUB_URL, verbose=0)

hubKey = client.hub.loginWithAuthRelayMode(HUB_LOGIN, HUB_PASSWORD)

#Get the server Ids
serverIds = client.hub.listServerIds(hubKey)

#authenticate those servers(same credentials will be used as of hub to authenticate)
client.hub.attachToServers(hubKey, serverIds)

# perform the needed operation 
systemsPerServer = client.multicast.system.list_systems(hubKey, serverIds)
successfulResponses = systemsPerServer["Successfull"]["Responses"]
failedResponses = systemsPerServer["Failed"]["Responses"]

for system in itertools.chain.from_iterable(successfulResponses):
  print (system)

#logout
client.auth.logout(hubKey)
