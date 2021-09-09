package overlord

import (
	"log"
	"os"
)

type State struct {
	Libvirt Libvirt
	Config  Config
}

func NewState() State {
	return State{
		Libvirt{},
		Config{},
	}
}

func (state *State) Setup() {
	state.Config.ReadConfig()
	if state.Config.LogFilePath != "" {
		file, err := os.OpenFile(state.Config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("Could not open log file! %s", state.Config.LogFilePath)
		} else {
			log.SetOutput(file)
		}
	}
	state.Libvirt.ConnectMany(state.Config.LibvirtHosts)
}

func (state *State) Stop() {
	for _, conn := range state.Libvirt.Hosts {
		conn.Close()
	}
}
