package main

import (
	"log"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/call"
)

func BringUpPiSwitch() {
	if !tracker.PiSwitchIsOff {
		log.Println("Pi Switch is already on!")
		return
	}
	log.Println("Bringing up Pi Switch")
	err := call.CallHassService("switch", "turn_on", "{\"entity_id\": \"switch.pi_switch_plug\"}")
	if err != nil {
		log.Printf("Error bringing up Pi Switch")
	} else {
		log.Println("Pi Switch bring up complete")
		tracker.PiSwitchIsOff = false
	}
}

func ShutdownPiSwitch() {
	if tracker.PiSwitchIsOff {
		log.Println("Pi Switch is already off!")
		return
	}
	log.Println("Shutting down Pi Switch")
	err := call.CallHassService("switch", "turn_off", "{\"entity_id\": \"switch.pi_switch_plug\"}")
	if err != nil {
		log.Printf("Error shutting down Pi Switch")
	} else {
		log.Println("Pi Switch shutdown complete")
		tracker.PiSwitchIsOff = true
	}
}