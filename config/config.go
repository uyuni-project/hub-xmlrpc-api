package config

import (
	"log"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

type HubConfig struct {
	SUMA_API_URL string
}

type Config struct {
	Hub                              HubConfig
	RelayMode                        bool
	AutoConnectMode                  bool
	ConnectTimeout, ReadWriteTimeout int
	//TODO:WE SHOULD GET THIS FROM SUMA API (ie, on listUserSystems)
	ServerURLByServerID map[string]string
}

// New returns a new Config struct
func New() *Config {
	if hubConfigFile, exists := os.LookupEnv("HUB_CONFIG_FILE"); exists {
		if err := k.Load(file.Provider(hubConfigFile), json.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		return &Config{
			Hub: HubConfig{
				SUMA_API_URL: k.String("hub.manager_api_url"),
			},
			RelayMode:        k.Bool("relay_mode"),
			AutoConnectMode:  k.Bool("auto_connect_mode"),
			ConnectTimeout:   k.Int("connect_timeout"),
			ReadWriteTimeout: k.Int("read_write_timeout"),
			//TODO: REMOVE IT. WE SHOULD GET THIS FROM SUMA API (ie, on listUserSystems)
			ServerURLByServerID: k.StringMap("endpoints"),
		}
	}
	return &Config{}
}
