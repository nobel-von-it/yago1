package main

import (
	"flag"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
	Same          bool
}

var config Config
var sugar *zap.SugaredLogger

func init() {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		defServAddr = serverAddress
	}
	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		defBaseUrl = baseUrl
	}

	flag.StringVar(&config.ServerAddress, "a", defServAddr, "address for http-server")
	flag.StringVar(&config.BaseUrl, "p", defBaseUrl, "user-input value for pre short url")
	flag.BoolVar(&config.Same, "s", false, "are base url and server address same")

	var logger *zap.Logger
	var err error
	if isTestRun() {
		logger = zap.NewNop()
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	sugar = logger.Sugar()

	if !isTestRun() {
		flag.Parse()
	}

	if !isTestRun() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}
}

func isTestRun() bool {
	for _, arg := range os.Args {
		if arg == "-test.v" || arg == "-test.run" {
			return true
		}
	}
	return false
}

var defServAddr = "localhost:8080"
var defBaseUrl = "localhost:8080"

func ParseArgs() *Config {
	if config.Same && config.ServerAddress != config.BaseUrl {
		if config.ServerAddress != defServAddr {
			config.BaseUrl = config.ServerAddress
		} else if config.BaseUrl != defBaseUrl {
			config.ServerAddress = config.BaseUrl
		}
	}

	return &config
}

func DefaultValues() {
	_ = flag.Set("a", defServAddr)
	_ = flag.Set("p", defBaseUrl)
}
