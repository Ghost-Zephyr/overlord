package main

import (
	log "bitsnthings.dev/overlord/src/log"
	overlord "bitsnthings.dev/overlord/src/state"
)

func main() {
	log.PrintLog(log.INFO, "Starting Overlord.")
	state := overlord.NewState()
	state.Setup()
	//state.MainLoop()
	state.Stop()
	log.PrintLog(log.TRACE, "Done.")
}
