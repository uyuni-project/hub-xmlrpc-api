#!/usr/bin/python
import sys
import xmlrpclib
import datetime    

api_url = "http://localhost:8000/hub/rpc/api
client = xmlrpclib.Server(api_url, verbose=0)

hubKey = client.hub.login("admin", "admin" )
print hubKey

serverIds = client.hub.listServerIds(hubKey)

print serverIds

usernames = ["admin" for s in serverIds]
passwords = ["admin" for s in serverIds]

client.Hub.attachToServers(hubKey, serverIds, usernames, passwords)

systemsPerServer = client.system.listUserSystems(hubKey, serverIds, [usernames])
print systemsPerServer

#schedule install a patch
package_ids = [[35, 36] for s in serverIds]
minions_per_server = client.system.listSystems(hubKey, serverIds)

minion_ids = [[minion["id"] for minion in minions] for minions in minions_per_server]
now = datetime.datetime.now()
start_times = [now + datetime.timedelta(days=360) for s in serverIds]

action_id = client.system.schedulePackageInstall(hubKey, serverIds, [minion_ids, package_ids, start_times])

print(action_id)
