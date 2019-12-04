#!/usr/bin/python
import sys
import xmlrpclib
import datetime	

api_url = "http://localhost:8000/RPC2"
client = xmlrpclib.Server(api_url, verbose=0)

loginResponse = client.Auth.Login("admin", "admin")

hubKey = loginResponse["hubSessionKey"]
print hubKey

serverKeys = loginResponse["serverSessionKeys"]
print serverKeys

serverKey1 = serverKeys[0]["sessionKey"]

systems = client.system.listUserSystems(hubKey, [[serverKey1, "admin"]])
print systems



