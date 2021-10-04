package state

import (
	"context"
	"fmt"
	logger "log"
	"os"
	"os/signal"
	"syscall"

	virt "bitsnthings.dev/overlord/src/libvirt"
	log "bitsnthings.dev/overlord/src/log"
	"bitsnthings.dev/overlord/src/matrix"
	"github.com/savsgio/atreugo/v11"
	"go.mongodb.org/mongo-driver/bson"
	"libvirt.org/go/libvirt"
)

func (state *State) Setup() {
	log.PrintLog(log.TRACE, "Staring setup.")
	state.setupSignalHandlers()
	state.Config.ReadConfig()
	state.setLogOutput()
	log.PrintLog(log.TRACE, "Config and logfile ready.")
	state.setupDB()
	log.PrintLog(log.INFO, "Connected to database.")
	state.Libvirt.ConnectMany(state.Config.LibvirtHosts, libvirt.NewConnect)
	state.Libvirt.ConnectMany(state.Config.LibvirtROHosts, libvirt.NewConnectReadOnly)
	log.PrintLog(log.INFO, "Done with inital libvirt connections.")
	state.Libvirt.GetAllDomains()
	log.PrintLog(log.TRACE, "Fetched cluster status and updated internal state.")
	state.updateDomains()
	log.PrintLog(log.TRACE, "Updated cluster state.")
	state.pushDBState()
	log.PrintLog(log.TRACE, "Pushed updated state to DB.")
	if state.Config.EnableAPI {
		state.API.Setup(&atreugo.Config{
			Addr: state.Config.APIBindAddress,
		})
		log.PrintLog(log.INFO, "Configured the HTTP API thread.")
	}
	if state.Config.EnableMatrix {
		matrix.Setup(state.Config)
		log.PrintLog(log.INFO, "Connected to matrix.")
	}
}

func (state *State) updateDomains() {
	var doms []virt.DomainDocument
	cur, err := state.DB.domains.Find(context.TODO(), bson.D{})
	if err != nil {
		log.PrintLog(log.ERROR, "Error tying to get all domains from database. ", err)
		return
	}
	err = cur.All(context.TODO(), &doms)
	if err != nil {
		log.PrintLog(log.ERROR, "Error decoding all domains after running the query. ", err)
	}
	for i := range doms {
		dom := doms[i]
		domain, active, err := state.Libvirt.GetDomainAndStatus(dom.URI, dom.UUID)
		if err == nil {
			if active {
				state.SetDomainState(libvirt.DOMAIN_RUNNING, domain)
			} else {
				state.SetDomainState(libvirt.DOMAIN_SHUTDOWN, domain)
			}
		}
	}
}

func (state *State) setupSignalHandlers() {
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
