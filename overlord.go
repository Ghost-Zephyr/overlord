package main

import (
	log "bitsnthings.dev/overlord/src/log"
	overlord "bitsnthings.dev/overlord/src/state"
)

func main() {
	log.PrintLog(log.INFO, "Starting Overlord.")
	state := overlord.NewState()
	state.Setup()
	log.PrintLog(log.INFO, "Setup done, entering main loops.")
	if state.Config.EnableAPI {
		go state.API.Start()
	}
	state.MainLoop()
	state.Stop()
	log.PrintLog(log.TRACE, "Done.")
}
