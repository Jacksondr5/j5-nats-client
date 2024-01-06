package logic

import (
	log "github.com/sirupsen/logrus"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logger"
)

type KillableDevice interface {
	devices.Device
	TurnOff() (err error)
}

func TurnOffDevice (device KillableDevice) {
	logFields := log.Fields{
		"deviceName": device.Name(),
	}
	if device.IsOff() {	
		logger.WarningWithFields("Device is already off", logFields)
		return
	}
	logger.InfoWithFields("Turning off device", logFields)
	err := device.TurnOff()
	if err != nil {
		logger.ErrorWithFields("Error turning off device", err, logFields)
	} else {
		logger.InfoWithFields("Successfully turned device off", logFields)
		device.SetIsOff(true)
	}
}