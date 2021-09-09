package state

import (
	"context"
	"fmt"
	logger "log"
	"os"
	"os/signal"
	"syscall"

	log "bitsnthings.dev/overlord/src/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (state *State) Setup() {
	log.PrintLog(log.TRACE, "Staring setup.")
	state.setupCloseSignalHandlers()
	state.Config.ReadConfig()
	state.setLogOutput()
	log.PrintLog(log.TRACE, "Config and logfile ready.")
	state.setupDB()
	log.PrintLog(log.TRACE, "Connected to database.")
	state.Libvirt.ConnectMany(state.Config.LibvirtHosts)
	log.PrintLog(log.TRACE, "Done with inital connections.")
	state.Libvirt.GetStatus()
	log.PrintLog(log.TRACE, "Fetched cluster status and prepared internal state.")
}

func (state *State) setupDB() {
	clientOptions := options.Client().ApplyURI(state.Config.MongoDbStr)
	client, _ := mongo.Connect(context.TODO(), clientOptions)
	err := client.Ping(context.TODO(), nil)
	// No idea why, but error only gets returned with ping, not connect
	// if the server refuses the connection.
	if err != nil {
		log.PrintLog(log.FATAL, "Error connecting to database wtih connection string \"%s\"! %s",
			state.Config.MongoDbStr, err)
	}
	state.Mongo = client
}

func (state *State) setupCloseSignalHandlers() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	signal.Notify(s, os.Interrupt, syscall.SIGINT)
	go func() {
		<-s
		fmt.Println()
		log.PrintLog(log.DEBUG, "Got interrupt!")
		state.Stop()
		os.Exit(0)
	}()
	log.PrintLog(log.TRACE, "Signal handlers ready.")
}

func (state *State) setLogOutput() {
	if state.Config.LogFilePath != "" {
		file, err := os.OpenFile(state.Config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.PrintLog(log.ERROR, "Could not open log file with path \"%s\"! %s",
				state.Config.LogFilePath, err)
		} else {
			logger.SetOutput(file)
		}
	} else {
		logger.SetOutput(os.Stdout)
	}
}
