package config

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")
var exampleJsonContent = `
{
	"type":"json",
	"hub":{
	   "manager_api_url":"http://localhost/rpc/api"
	},
	"connect_timeout":10,
	"read_write_timeout":10
 }
`

type HubConfig struct {
	SUMA_API_URL string
}

type Config struct {
	Hub                              HubConfig
	ConnectTimeout, ReadWriteTimeout int
}

// InitializeConfig returns a new Config struct
func InitConfig() *Config {
	if hubConfigFile, exists := os.LookupEnv("HUB_CONFIG_FILE"); exists {
		if err := k.Load(file.Provider(hubConfigFile), json.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		return &Config{
			Hub: HubConfig{
				SUMA_API_URL: k.String("hub.manager_api_url"),
			},
			ConnectTimeout:   k.Int("connect_timeout"),
			ReadWriteTimeout: k.Int("read_write_timeout"),
		}
	}
	errorToPrint := fmt.Errorf("Please set the HUB_CONFIG_FILE env variable which point to a json file."+
		" An example content is given below\n %s", exampleJsonContent)
	panic(errorToPrint)
}
