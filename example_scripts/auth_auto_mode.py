#!/usr/bin/python
import xmlrpclib  
import itertools

HUB_URL = "http://localhost:8888/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"

client = xmlrpclib.Server(HUB_URL, verbose=0)

hubKey = client.hub.loginWithAutoconnectMode(HUB_LOGIN, HUB_PASSWORD)

#Get the server Ids
serverIds = client.hub.listServerIds(hubKey)

# perform the needed operation 
systemsPerServer = client.multicast.system.list_systems(hubKey, serverIds)
successfulResponses = systemsPerServer["Successfull"]["Responses"]
failedResponses = systemsPerServer["Failed"]["Responses"]

for system in itertools.chain.from_iterable(successfulResponses):
  print (system)

#logout
client.auth.logout(hubKey)