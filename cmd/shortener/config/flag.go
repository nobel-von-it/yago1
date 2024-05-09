package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
	Same          bool
}

var config = new(Config)

var defServAddr = "localhost:8080"
var defBaseUrl = "localhost:8080"

func ParseArgs() *Config {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		defServAddr = serverAddress
	}
	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		defBaseUrl = baseUrl
	}

	flag.StringVar(&config.ServerAddress, "a", defServAddr, "address for http-server")
	flag.StringVar(&config.BaseUrl, "p", defBaseUrl, "user-input value for pre short url")
	flag.BoolVar(&config.Same, "s", false, "are base url and server address same")

	flag.Parse()

	if config.Same && config.ServerAddress != config.BaseUrl {
		if config.ServerAddress != defServAddr {
			config.BaseUrl = config.ServerAddress
		} else if config.BaseUrl != defBaseUrl {
			config.ServerAddress = config.BaseUrl
		}
	}

	return config
}
