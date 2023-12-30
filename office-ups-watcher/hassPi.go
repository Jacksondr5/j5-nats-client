package main

import "log"

func ShutdownHass() {
	if tracker.HassPiIsOff {
		log.Println("Home Assistant Pi is already off!")
		return
	}
	log.Println("Shutting down Home Assistant Pi")
	err := CallHassService("hassio", "host_shutdown", "")
	if err != nil {
		log.Printf("Error shutting down Home Assistant Pi")
	} else {
		log.Println("Home Assistant Pi shutdown complete")
		tracker.HassPiIsOff = true
	}
}