package main

import (
	"log"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/call"
)

func ShutdownHass() {
	if tracker.HassPiIsOff {
		log.Println("Home Assistant Pi is already off!")
		return
	}
	log.Println("Shutting down Home Assistant Pi")
	err := call.CallHassService("hassio", "host_shutdown", "")
	if err != nil {
		log.Printf("Error shutting down Home Assistant Pi")
	} else {
		log.Println("Home Assistant Pi shutdown complete")
		tracker.HassPiIsOff = true
	}
}