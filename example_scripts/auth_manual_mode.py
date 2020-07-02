#!/usr/bin/python
import xmlrpclib  
import itertools

HUB_URL = "http://localhost:2830/hub/rpc/api"
HUB_LOGIN = "admin"
HUB_PASSWORD = "admin"

client = xmlrpclib.Server(HUB_URL, verbose=0)

hubSessionKey = client.hub.login(HUB_LOGIN, HUB_PASSWORD)

#get the server IDs
serverIDs = client.hub.listServerIds(hubSessionKey)

#authenticate those servers(each server with it's own credentials)
usernames = ["admin" for s in serverIDs]
passwords = ["admin" for s in serverIDs]
client.hub.attachToServers(hubSessionKey, serverIDs, usernames, passwords)

#perform the required operation 
systemsPerServer = client.multicast.system.listSystems(hubSessionKey, serverIDs)
successfulResponses = systemsPerServer["Successful"]["Responses"]
failedResponses = systemsPerServer["Failed"]["Responses"]

for system in itertools.chain.from_iterable(successfulResponses):
    print(system)

#logout
client.hub.logout(hubSessionKey)
