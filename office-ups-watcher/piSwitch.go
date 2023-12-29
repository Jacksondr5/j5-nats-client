package main

import "log"

func BringUpPiSwitch() {
	if !tracker.PiSwitchIsOff {
		log.Println("PiSwitch is already on!")
		return
	}
	tracker.PiSwitchIsOff = false
	panic("unimplemented")
}

func ShutdownPiSwitch() {
	if tracker.PiSwitchIsOff {
		log.Println("PiSwitch is already off!")
		return
	}
	tracker.PiSwitchIsOff = true
	panic("unimplemented")
}