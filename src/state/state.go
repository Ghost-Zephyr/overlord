package state

import (
	"context"
	"time"

	virt "bitsnthings.dev/overlord/src/libvirt"
	log "bitsnthings.dev/overlord/src/log"
	"bitsnthings.dev/overlord/src/state/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Version = "0.0.69"

type State struct {
	Libvirt       virt.Libvirt
	Config        conf.Config
	MongoDB       *mongo.Database
	upsertOpts    *options.UpdateOptions
	streamOpts    *options.ChangeStreamOptions
	changeStreams []*mongo.ChangeStream
}

func NewState() State {
	return State{
		virt.New(),
		conf.Config{},
		&mongo.Database{},
		&options.UpdateOptions{},
		&options.ChangeStreamOptions{},
		[]*mongo.ChangeStream{},
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
	for _, stream := range state.changeStreams {
		stream.Close(context.TODO())
	}
	log.PrintLog(log.TRACE, "Closed mongodb change streams.")
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
}
