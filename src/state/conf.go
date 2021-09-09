package state

import (
	"encoding/json"
	"flag"
	"os"

	log "bitsnthings.dev/overlord/src/log"
)

type Config struct {
	LibvirtHosts []string
	MongoDbStr   string
	ConfFilePath string
	LogLevel     log.LogLevel
	LogFilePath  string
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
