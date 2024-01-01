package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/call"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/nats-io/nats.go"
)



func main() {
	log.Println("Starting UPS Watcher")
	hostname := getHostname()
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222", nats.Name(hostname))
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
	k8sPiCountChan := make(chan int)
	nc.Subscribe("ups.office.ack", func(m *nats.Msg) {
		k8sPiCount = logic.OnShutdownAck(m, &devices, k8sPiCount)
		k8sPiCountChan <- k8sPiCount
	})

	log.Println("Subscription setup complete.  Polling battery status.")

	go poll(&tracker, nc, &devices)
	http.ListenAndServe(":12346", nil)
}

func poll(tracker *logic.Tracker, nc *nats.Conn, devices *logic.ManagedDevices) {
	for {
		sleepTime := 10 * time.Second
			sleepTime = logic.ExecutePollingLogic(tracker, nc, devices, battery.BatteryPollerImpl{})
			// I think theres a race condition here with the subscription and polling logic
			if tracker.BadBatteryStatusCount > 5 {
				log.Fatalln("Too many errors polling battery status, exiting.")
			}
		time.Sleep(sleepTime)
	}
}

func getHostname() string {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		panic(hostnameErr)
	}
	return hostname
}
