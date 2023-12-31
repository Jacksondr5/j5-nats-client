package main

import (
	"log"
	"time"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/call"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/nats-io/nats.go"
)



func main() {
	log.Println("Starting UPS Watcher")
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222")
	defer nc.Drain()
	log.Println("Connected to NATS")

	log.Println("Setting up subscriptions")
	var tracker = logic.Tracker{
		BadBatteryStatusCount: 0,
		Group1IsDeactivated: false,
		Group2IsDeactivated: false,
		Group3IsDeactivated: false,
		IsActive: false,
	}
	devices := logic.ManagedDevices{
		GitLabPi: devices.NewUnmanageableDevice("GitLab Pi"),
		Nas: devices.NewApiManagedDevice(
			"NAS", 
			"5-bZm7aVBKJbTGszxNBRolUDnzqezCRG3R83S27L6ztc5rTcpg6JKG01OPbjDRnjjq", 
			"http://nas.j5/api/v2.0/system/shutdown", 
			"", 
			"",
			call.HttpClientImpl{},
		),
		PiSwitch: devices.NewApiManagedDevice(
			"Pi Switch", 
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlNjkyMGRhMDc5NWU0ZThmOGUzYzYyOTAzYzgwZmE0NyIsImlhdCI6MTcwMzg4MjQ0OCwiZXhwIjoyMDE5MjQyNDQ4fQ.BWphYONMeYF2Z64N6uAhhqNIOG3D8FfE3RjSR9XgrtM", 
			"http://hass.j5:8123/api/services/switch/turn_off", 
			"http://hass.j5:8123/api/services/switch/turn_on", 
			"{\"entity_id\": \"switch.pi_switch\"}",
			call.HttpClientImpl{},
		),
	}
	k8sPiCount := 0
	nc.Subscribe("ups.office.ack", func(m *nats.Msg) {
		k8sPiCount = logic.OnShutdownAck(m, &devices, k8sPiCount)
	})

	log.Println("Subscription setup complete.  Polling battery status.")

	for {
		sleepTime, newK8sPiCount := logic.ExecutePollingLogic(&tracker, k8sPiCount, nc, &devices, battery.BatteryPollerImpl{})
		// I think theres a race condition here with the subscription and polling logic
		if newK8sPiCount == 0 {
			k8sPiCount = newK8sPiCount
		}
		if tracker.BadBatteryStatusCount > 5 {
			log.Fatalln("Too many errors polling battery status, exiting.")
		}
		time.Sleep(sleepTime)
	}
}


