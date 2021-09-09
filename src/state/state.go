package state

import (
	virt "bitsnthings.dev/overlord/src/libvirt"
	log "bitsnthings.dev/overlord/src/log"
	"go.mongodb.org/mongo-driver/mongo"
)

const Version = "0.0.69"

type State struct {
	Libvirt virt.Libvirt
	Config  Config
	Mongo   *mongo.Client
}

func NewState() State {
	return State{
		virt.New(),
		Config{},
		&mongo.Client{},
	}
}

func (state *State) MainLoop() {
	state.Libvirt.GetStatus()
}

func (state *State) Stop() {
	log.PrintLog(log.INFO, "Stopping Overlord.")
	for _, conn := range state.Libvirt.Hosts {
		conn.Close()
	}
	log.PrintLog(log.TRACE, "All libvirt connections closed.")
}
