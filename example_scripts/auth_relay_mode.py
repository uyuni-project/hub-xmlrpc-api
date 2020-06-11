#!/usr/bin/python
import xmlrpclib  
import itertools

HUB_URL = "http://localhost:2830/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"

client = xmlrpclib.Server(HUB_URL, verbose=0)

hubSessionKey = client.hub.loginWithAuthRelayMode(HUB_LOGIN, HUB_PASSWORD)

#get the server IDs
serverIDs = client.hub.listServerIds(hubSessionKey)

#authenticate those servers(same credentials will be used as of hub to authenticate)
client.hub.attachToServers(hubSessionKey, serverIDs)

#perform the required operation 
systemsPerServer = client.multicast.system.listSystems(hubSessionKey, serverIDs)
successfulResponses = systemsPerServer["Successful"]["Responses"]
failedResponses = systemsPerServer["Failed"]["Responses"]

for system in itertools.chain.from_iterable(successfulResponses):
  print(system)

#logout
client.auth.logout(hubSessionKey)
