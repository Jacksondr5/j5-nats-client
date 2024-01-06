package logic

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logger"
)

type Tracker struct {
	BadBatteryStatusCount int
	Group1IsDeactivated bool
	Group2IsDeactivated bool
	Group3IsDeactivated bool
	IsActive bool
}

type ManagableDevice interface {
	KillableDevice
	RaisableDevice
}

type ManagedDevices struct {
	GitLabPi devices.Device
	Nas KillableDevice
	PiSwitch ManagableDevice
}

type NatsClient interface {
	Publish(subject string, data []byte) error
}

type BatteryPoller interface {
	PollBatteryStatus() (battery.BatteryStatus, error)
}

func ExecutePollingLogic(tracker *Tracker, nc NatsClient, devices *ManagedDevices, batteryPoller BatteryPoller) (time.Duration) {
	batteryStatus, err := batteryPoller.PollBatteryStatus();
	sleepTime := 10 * time.Second
	if err != nil {
		logger.Error("Error polling battery status", err)
		tracker.BadBatteryStatusCount++
		return sleepTime
	}
	logger.DebugWithFields("Battery status", log.Fields{
		"isOnBattery": batteryStatus.IsOnBattery,
		"percent": batteryStatus.Percent,
	})
	if batteryStatus.IsOnBattery && batteryStatus.Percent <= 95 {
		logBatteryStatus("System is on battery.  Turning things off.", batteryStatus)
		sleepTime = 1 * time.Second
		tracker.IsActive = true
		if !tracker.Group1IsDeactivated {
			deactivateViaNats(nc, "ups.office", "group1", &tracker.Group1IsDeactivated)
		}
		if !tracker.Group2IsDeactivated && batteryStatus.Percent <= 85 {
			deactivateViaNats(nc, "ups.office", "group2", &tracker.Group2IsDeactivated)
		}
		if !tracker.Group3IsDeactivated && batteryStatus.Percent <= 40 {
			logger.Info("Deactivating group 3")
			TurnOffDevice(devices.PiSwitch)
			TurnOffDevice(devices.Nas)
			tracker.Group3IsDeactivated = true
		}
	} else if tracker.IsActive && !batteryStatus.IsOnBattery && batteryStatus.Percent >= 95 {
		logBatteryStatus("Battery is charged.  Turning things back on.", batteryStatus)
		tracker.Group1IsDeactivated = false
		tracker.Group2IsDeactivated = false
		tracker.Group3IsDeactivated = false
		if !devices.Nas.IsOff() {
			logger.Info("NAS is not off, turning its dependents back on")
			TurnOnDevice(devices.PiSwitch)
		} else {
			logger.Info("NAS is off, not turning its dependents back on")
		}
		tracker.IsActive = false
	} else if tracker.IsActive && !batteryStatus.IsOnBattery {
		logBatteryStatus("System is on line power, but still charging battery.  Waiting for battery to charge.", batteryStatus)
	} else {
		logger.Debug("Everything is normal")
	}
	return sleepTime
}

func deactivateViaNats(nc NatsClient, subject string, data string, trackerFlag *bool) {
	logger.Info(fmt.Sprintf("Deactivating group %s", data))
	err := nc.Publish(subject, []byte(data))
	if err != nil {
		if err.Error() == "nats: invalid connection" {
			logger.Fatal("Connection to nats is invalid, exiting", err)
		}
		logger.Error("Error publishing to nats", err)
		return
	}
	*trackerFlag = true
}

func logBatteryStatus(msg string, batteryStatus battery.BatteryStatus) {
	logger.DebugWithFields(msg, log.Fields{
		"isOnBattery": batteryStatus.IsOnBattery,
		"percent": batteryStatus.Percent,
	})
}
