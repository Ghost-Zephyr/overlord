package overlord

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Config struct {
	LibvirtHosts []string
	PrivateKey   string
	ConfFilePath string
	LogFilePath  string
}

func (conf *Config) ReadConfig() {
	flag.StringVar(&conf.ConfFilePath, "c", "lord.json", "Config file path.")
	flag.StringVar(&conf.LogFilePath, "l", "", "Log file path.")
	flag.Parse()
	file, err := os.Open(conf.ConfFilePath)
	if err != nil {
		log.Printf("Error opening config file! %s", err)
	} else {
		defer file.Close()
		json.NewDecoder(file).Decode(&conf)
	}
}
