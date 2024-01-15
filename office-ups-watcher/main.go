package main

import (
	"net/http"
	"os"
	"time"

	httpclient "github.com/jacksondr5/go-monorepo/httpclient"
	"github.com/jacksondr5/go-monorepo/logger"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/nats-io/nats.go"
)



func main() {
	logger.Init()
	logger.Info("Starting UPS Watcher")
	hostname := getHostname()
	nc, err := nats.Connect(
		"nats://nats.k8s.j5:4222", 
		nats.Name(hostname), 
		nats.MaxReconnects(-1),
	)
	if err != nil {
		logger.Fatal("Error connecting to NATS", err)
	}
	defer nc.Drain()
	logger.Info("Connected to NATS")

	logger.Info("Setting up subscriptions")
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
			httpclient.HttpClientImpl{},
		),
		PiSwitch: devices.NewApiManagedDevice(
			"Pi Switch", 
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlNjkyMGRhMDc5NWU0ZThmOGUzYzYyOTAzYzgwZmE0NyIsImlhdCI6MTcwMzg4MjQ0OCwiZXhwIjoyMDE5MjQyNDQ4fQ.BWphYONMeYF2Z64N6uAhhqNIOG3D8FfE3RjSR9XgrtM", 
			"http://hass.j5:8123/api/services/switch/turn_off", 
			"http://hass.j5:8123/api/services/switch/turn_on", 
			"{\"entity_id\": \"switch.pi_switch_plug\"}",
			httpclient.HttpClientImpl{},
		),
	}

	actionFeed := make(chan logic.PiMonitorAction)
	monitorChannel := make(chan int)
	go logic.MonitorPiShutdown(actionFeed, monitorChannel, &devices, 12)
	_, err = nc.Subscribe("ups.office.ack", func(m *nats.Msg) {
		logic.OnShutdownAck(m, actionFeed, monitorChannel)
	})
	if err != nil {
		logger.Fatal("Error setting up subscription for \"ups.office.ack\"", err)
	}

	logger.Info("Subscription setup complete.  Polling battery status.")
	go poll(&tracker, nc, &devices, actionFeed)
	http.ListenAndServe(":12346", nil)
}

func poll(tracker *logic.Tracker, nc *nats.Conn, devices *logic.ManagedDevices, actionFeed chan<- logic.PiMonitorAction) {
	for {
		sleepTime := 10 * time.Second
			sleepTime = logic.ExecutePollingLogic(tracker, nc, devices, battery.BatteryPollerImpl{}, actionFeed)
			// I think theres a race condition here with the subscription and polling logic
			if tracker.BadBatteryStatusCount > 5 {
				logger.Fatal("Too many errors polling battery status, exiting.", nil)
			}
		time.Sleep(sleepTime)
	}
}

func getHostname() string {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		logger.Fatal("Error getting hostname", hostnameErr)
	}
	return hostname
}
