package main

import "log"


func ShutdownNas() {
	const accessToken = "5-bZm7aVBKJbTGszxNBRolUDnzqezCRG3R83S27L6ztc5rTcpg6JKG01OPbjDRnjjq"
	if tracker.NasIsOff {
		log.Println("NAS is already off!")
		return
	}
	log.Println("Shutting down NAS")
	err := HttpPost(
		"http://nas.j5/api/v2.0/system/shutdown",
		"",
		"NAS shutdown",
		accessToken,
	)
	if err != nil {
		log.Printf("Error shutting down NAS")
	} else {
		log.Println("NAS shutdown complete")
		tracker.NasIsOff = true
	}
}