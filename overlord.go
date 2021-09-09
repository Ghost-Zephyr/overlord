package main

import (
	"log"

	overlord "bitsnthings.dev/overlord/src"
)

func main() {
	log.Println("Starting Overlord.")
	state := overlord.NewState()
	state.Setup()
	log.Println("Stopping Overlord.")
	state.Stop()
}
