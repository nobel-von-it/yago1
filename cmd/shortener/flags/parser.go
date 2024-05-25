package flags

import (
	"flag"
	"go.uber.org/zap"
	"nerd/shortener/files"
	"nerd/shortener/storage"
	"os"
	"strings"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
	StoragePath   string
	Same          bool
}

var (
	defServAddr    = "localhost:8080"
	defBaseUrl     = "localhost:8080"
	defStoragePath = "tmp/short-url-db.json"
	config         Config
	sugar          *zap.SugaredLogger
	events         = storage.GetEvents()
)

func init() {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		defServAddr = serverAddress
	}
	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		defBaseUrl = baseUrl
	}
	if storagePath := os.Getenv("STORAGE_PATH"); storagePath != "" {
		defStoragePath = storagePath
	}

	flag.StringVar(&config.ServerAddress, "a", defServAddr, "address for http-server")
	flag.StringVar(&config.BaseUrl, "p", defBaseUrl, "user-input value for pre short url")
	flag.StringVar(&config.StoragePath, "f", defStoragePath, "full path to the storage")
	flag.BoolVar(&config.Same, "s", false, "are base url and server address same")

	if !isTestRun() {
		flag.Parse()
	}

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
	if err = events.Load(config.StoragePath); err != nil {
		panic(err)
	}

	if !isTestRun() {
		err := logger.Sync()
		if err != nil && !strings.Contains(err.Error(), "stderr") {
			panic(err)
		}
	}

	if config.Same && config.ServerAddress != config.BaseUrl {
		if config.ServerAddress != defServAddr {
			config.BaseUrl = config.ServerAddress
		} else if config.BaseUrl != defBaseUrl {
			config.ServerAddress = config.BaseUrl
		}
	}

	if err = files.DirExists(config.StoragePath); err != nil {
		panic(err)
	}
}

func isTestRun() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "test") {
			return true
		}
	}
	return false
}

func GetConfig() Config {
	return config
}

func GetSugar() *zap.SugaredLogger {
	return sugar
}
