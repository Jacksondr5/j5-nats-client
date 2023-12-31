package logic

import (
	"log"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
)

type RaisableDevice interface {
	devices.Device
	TurnOn() (err error)
}

func TurnOnDevice (device RaisableDevice) {
	if !device.IsOff() {	
		log.Printf("%s is already on!", device.Name())
		return
	}
	log.Printf("Bringing up %s", device.Name())
	err := device.TurnOn()
	if err != nil {
		log.Printf("Error bringing up %s", device.Name())
	} else {
		log.Printf("%s bring up complete", device.Name())
		device.SetIsOff(false)
	}
}