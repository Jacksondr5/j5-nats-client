package logic

import (
	log "github.com/sirupsen/logrus"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logger"
)

type RaisableDevice interface {
	devices.Device
	TurnOn() (err error)
}

func TurnOnDevice (device RaisableDevice) {
	logFields := log.Fields{
		"deviceName": device.Name(),
	}
	if !device.IsOff() {	
		logger.WarningWithFields("Device is already on", logFields)
		return
	}
	logger.InfoWithFields("Turning on device", logFields)
	err := device.TurnOn()
	if err != nil {
		logger.ErrorWithFields("Error turning on device", err, logFields)
	} else {
		logger.InfoWithFields("Successfully turned device on", logFields)
		device.SetIsOff(false)
	}
}