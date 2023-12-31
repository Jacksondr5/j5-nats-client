package logic

import (
	"log"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
)

type KillableDevice interface {
	devices.Device
	TurnOff() (err error)
}

func TurnOffDevice (device KillableDevice) {
	if device.IsOff() {	
		log.Printf("%s is already off!", device.Name())
		return
	}
	log.Printf("Shutting down %s", device.Name())
	err := device.TurnOff()
	if err != nil {
		log.Printf("Error shutting down %s", device.Name())
	} else {
		log.Printf("%s shutdown complete", device.Name())
		device.SetIsOff(true)
	}
}