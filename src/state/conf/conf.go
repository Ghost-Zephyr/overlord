package conf

import (
	"encoding/json"
	"flag"
	"os"

	log "bitsnthings.dev/overlord/src/log"
	"maunium.net/go/mautrix"
)

type Config struct {
	LibvirtHosts         []string
	LibvirtReadOnlyHosts []string
	MongoDbStr           string
	MongoDbName          string
	ConfFilePath         string
	LogLevel             log.LogLevel
	LogFilePath          string
	EnableMatrix         bool
	MatrixCreds          MatrixCreds
	MatrixCleint         mautrix.Client
}

type MatrixCreds struct {
	Homeserver string
	// AccessToken string
	// UserID      mautrix.UserIdentifier
	Username string
	Password string
}

func (conf *Config) ReadConfig() {
	flag.StringVar(&conf.ConfFilePath, "c", "lord.json", "Config file path.")
	flag.StringVar(&conf.LogFilePath, "l", "", "Log file path.")
	flag.Parse()
	file, err := os.Open(conf.ConfFilePath)
	if err != nil {
		log.PrintLog(log.ERROR, "Error opening config file! %s", err)
	} else {
		defer file.Close()
		json.NewDecoder(file).Decode(&conf)
	}
	log.CurrentLogLevel = conf.LogLevel
	log.PrintLog(log.DEBUG, "Log level: %s", log.CurrentLogLevel.String())
}
