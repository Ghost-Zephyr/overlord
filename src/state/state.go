/*
Package state for the overlord software.
This is where much of the magic happens!
*/
package state

import (
	"context"
	"time"

	"bitsnthings.dev/overlord/src/api"
	virt "bitsnthings.dev/overlord/src/libvirt"
	log "bitsnthings.dev/overlord/src/log"
	"bitsnthings.dev/overlord/src/state/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewState() State {
	return State{
		virt.New(),
		conf.Config{},
		DB{},
		&options.UpdateOptions{},
		&options.ChangeStreamOptions{},
		[]*mongo.ChangeStream{},
		api.API{},
	}
}

func (state *State) MainLoop() {
	for {
		time.Sleep(time.Second * 10)
	}
	//	state.Libvirt.GetStatus()
	//	state.Config.MatrixCleint.Sync()
}

func (state *State) Stop() {
	log.PrintLog(log.INFO, "Stopping Overlord.")
	// Close libvritd connections
	for uri, domMap := range state.Libvirt.Domains.Active {
		log.PrintLog(
			log.TRACE,
			"Freeing internal domain structures for libvirt domain with uri: %s.",
			uri)
		for _, dom := range domMap {
			dom.Free()
		}
	}
	for _, conn := range state.Libvirt.Hosts {
		conn.Close()
	}
	log.PrintLog(log.TRACE, "All libvirt connections closed.")
	// MongoDB clean up
	if !state.Config.InMemoryDB {
		for _, stream := range state.changeStreams {
			stream.Close(context.TODO())
		}
		log.PrintLog(log.TRACE, "Closed mongodb change streams.")
	}
}
