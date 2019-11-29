package main

import (
	"flag"
	"log"
	"os"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/gorilla"
	"github.com/devopsfaith/krakend/router/mux"
)

func main() {
	port := flag.Int("p", 8000, "Port of the service")
	logLevel := flag.String("l", "ERROR", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/home/marcelo/go/src/github.com/chiaradiamarcelo/hub_xmlrpc_api/krakend_conf.json", "Path to the configuration filename")
	flag.Parse()

	parser := config.NewParser()
	config.RoutingPattern = config.BracketsRouterPatternBuilder
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	/*
		secureMiddleware := secure.New(secure.Options{
			AllowedHosts:          []string{"127.0.0.1:8000", "example.com", "ssl.example.com"},
			SSLRedirect:           false,
			SSLHost:               "ssl.example.com",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			STSPreload:            true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
		})*/

	cfg := gorilla.DefaultConfig(customProxyFactory{logger, proxy.DefaultFactory(logger)}, logger)
	//cfg.Middlewares = append(cfg.Middlewares, secureMiddleware)
	routerFactory := mux.NewFactory(cfg)
	routerFactory.New().Run(serviceConfig)
}

// customProxyFactory adds a logging middleware wrapping the internal factory
type customProxyFactory struct {
	logger  logging.Logger
	factory proxy.Factory
}

// New implements the Factory interface
func (cf customProxyFactory) New(cfg *config.EndpointConfig) (p proxy.Proxy, err error) {
	p, err = cf.factory.New(cfg)
	if err == nil {
		p = proxy.NewLoggingMiddleware(cf.logger, cfg.Endpoint)(p)
	}
	return
}
