package logic

import (
	"log"
	"time"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
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

func ExecuteLogic(tracker *Tracker, nc NatsClient, devices *ManagedDevices, batteryPoller BatteryPoller) (time.Duration) {
	batteryStatus, err := batteryPoller.PollBatteryStatus();
	sleepTime := 10 * time.Second
	if err != nil {
		log.Println(err)
		log.Println("Error polling battery status, trying again....")
		tracker.BadBatteryStatusCount++
		return sleepTime
	}
	log.Printf("Battery status: %#v", batteryStatus)
	if batteryStatus.IsOnBattery && batteryStatus.Percent <= 95 {
		log.Println("System is on battery.  Turning things off.")
		sleepTime = 1 * time.Second
		tracker.IsActive = true
		if !tracker.Group1IsDeactivated {
			deactivateViaNats(nc, "ups.office", "group1", &tracker.Group1IsDeactivated)
		}
		if !tracker.Group2IsDeactivated && batteryStatus.Percent <= 85 {
			deactivateViaNats(nc, "ups.office", "group2", &tracker.Group2IsDeactivated)
		}
		if !tracker.Group3IsDeactivated && batteryStatus.Percent <= 40 {
			log.Println("Deactivating group 3")
			TurnOffDevice(devices.PiSwitch)
			TurnOffDevice(devices.Nas)
			tracker.Group3IsDeactivated = true
		}
	} else if tracker.IsActive && !batteryStatus.IsOnBattery && batteryStatus.Percent >= 95 {
		// Exit
		log.Println("Battery is charged.  Turning things back on.")
		if !devices.Nas.IsOff() {
			log.Println("NAS is not off, turning its dependents back on")
			TurnOnDevice(devices.PiSwitch)
		} else {
			log.Println("NAS is off, not turning its dependents back on")
		}
		tracker.IsActive = false
	} else if tracker.IsActive && !batteryStatus.IsOnBattery {
		log.Println("System is on line power, but still charging battery.  Waiting for battery to charge.")
	} else {
		log.Println("Everything is normal")
	}
	return sleepTime
}

func deactivateViaNats(nc NatsClient, subject string, data string, trackerFlag *bool) {
	log.Printf("Deactivating group %s", data)
	err := nc.Publish(subject, []byte(data))
	if err != nil {
		log.Println("Error publishing to nats")
		log.Println(err)
		return
	}
	*trackerFlag = true
}
