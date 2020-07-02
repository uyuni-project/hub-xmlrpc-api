package integration_tests

import (
	"os"

	"github.com/uyuni-project/hub-xmlrpc-api/initialization"
)

var (
	//Hub API Gateway server URL
	gatewayServerURL = "http://localhost:2830/hub/rpc/api"

	//1st peripheral server with 2 minions
	minion1ForServer1 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010000-minion-1",
	}
	minion2ForServer1 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010000-minion-2",
	}
	peripheralServer1 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010000",
		fqdn: "localhost:8002",
		minions: map[int64]SystemInfo{
			minion1ForServer1.id: minion1ForServer1,
			minion2ForServer1.id: minion2ForServer1,
		},
		port: 8002,
	}

	//2nd peripheral server with 2 minions
	minion1ForServer2 = SystemInfo{
		id:   1000010000,
		name: "peripheral-server-1000010001-minion-1",
	}
	minion2ForServer2 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010001-minion-2",
	}
	peripheralServer2 = SystemInfo{
		id:   1000010001,
		name: "peripheral-server-1000010001",
		fqdn: "localhost:8003",
		minions: map[int64]SystemInfo{
			minion1ForServer2.id: minion1ForServer2,
			minion2ForServer2.id: minion2ForServer2,
		},
		port: 8003,
	}

	peripheralServers = map[int64]SystemInfo{
		peripheralServer1.id: peripheralServer1,
		peripheralServer2.id: peripheralServer2,
	}
)

func init() {
	os.Setenv("HUB_API_URL", "http://localhost:8001/rpc/api")
	go func() {
		initialization.InitServer()
	}()
	initInfrastructure(peripheralServers, 8001, "admin", "admin")
}
