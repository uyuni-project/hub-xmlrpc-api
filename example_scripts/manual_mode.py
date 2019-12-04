#!/usr/bin/python
import sys
import xmlrpclib
import datetime	

api_url = "http://localhost:8000/RPC2"
client = xmlrpclib.Server(api_url, verbose=0)

loginResponse = client.Auth.Login("admin", "admin")
hubKey = loginResponse["hubSessionKey"]
print hubKey

serverKey1 = client.Auth.AttachToServer(hubKey, 1000010000, "admin", "admin")
print serverKey1

serverKey2 = client.Auth.AttachToServer(hubKey, 1000010001, "admin", "admin")
print serverKey2

systems = client.system.listUserSystems(hubKey, [[serverKey1, "admin"], [serverKey2, "admin"]])
print systems



