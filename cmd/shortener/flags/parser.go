package flags

import (
	"flag"
	"nerd/shortener/files"
	"nerd/shortener/storage"
	"os"
	"strings"

	"go.uber.org/zap"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
	StoragePath   string
	DataBase      string
	Same          bool
}

var (
	defServAddr    = "localhost:8080"
	defBaseUrl     = "localhost:8080"
	defStoragePath = "tmp/short-url-db.json"
	defDataBase    = "host=localhost port=5432 user=nerd password=123 dbname=shortener sslmode=disable"
	config         Config
	sugar          *zap.SugaredLogger
	events         = storage.GetEvents()
)

func checkenv(env, def string) string {
	if got := os.Getenv(env); got != "" {
		return got
	}
	return def
}

func init() {
	serverAddress := checkenv("SERVER_ADDRESS", defServAddr)
	baseUrl := checkenv("BASE_URL", defBaseUrl)
	storagePath := checkenv("STORAGE_PATH", defStoragePath)
	dataBase := checkenv("DATABASE_DSN", defDataBase)

	flag.StringVar(&config.ServerAddress, "a", serverAddress, "address for http-server")
	flag.StringVar(&config.BaseUrl, "p", baseUrl, "user-input value for pre short url")
	flag.StringVar(&config.StoragePath, "f", storagePath, "full path to the storage")
	flag.StringVar(&config.DataBase, "d", dataBase, "config for connecting to database")
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
