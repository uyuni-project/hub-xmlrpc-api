#!/usr/bin/python
import sys
import xmlrpclib
import datetime    

api_url = "http://localhost:8000/hub/rpc/api"
client = xmlrpclib.Server(api_url, verbose=0)

hubKey = client.hub.login("admin", "admin" )
print hubKey

serverIds = client.hub.listServerIds(hubKey)
print serverIds

usernames = ["admin" for s in serverIds]

systemsPerServer = client.multicast.system.listUserSystems(hubKey, serverIds, usernames)
print systemsPerServer
