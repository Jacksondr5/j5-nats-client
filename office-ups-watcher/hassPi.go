package main

import "log"

func BringUpHass() {
	log.Println("Turning HassPi back on")
	tracker.HassPiIsOff = false
	panic("unimplemented")
}

func ShutdownHass() {
	tracker.HassPiIsOff = true
	panic("unimplemented")
}