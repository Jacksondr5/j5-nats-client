package main

import (
	"log"
	"time"

	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/nats-io/nats.go"
)

type BatteryStatus struct {
	IsOnBattery bool
	Percent int
}

type Tracker struct {
	GitLabPiIsOff bool
	Group1Deactivated bool
	Group2Deactivated bool
	Group3Deactivated bool
	Group4Deactivated bool
	HassPiIsOff bool
	K8sPiCount int
	NasIsOff bool
	PiSwitchIsOff bool
	UdmIsOff bool
}
const totalK8sPis = 12

var tracker = Tracker{
	K8sPiCount: 0,
	GitLabPiIsOff: false,
	Group1Deactivated: false,
	Group2Deactivated: false,
	Group3Deactivated: false,
	Group4Deactivated: false,
	HassPiIsOff: false,
	UdmIsOff: false,
	PiSwitchIsOff: false,
}

func main() {
	//TODO: Figure out where this exe logs to
	log.Println("Starting UPS Watcher")
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222")
	defer nc.Drain()
	log.Println("Connected to NATS")

	log.Println("Setting up subscriptions")
	nc.Subscribe("ups.office.ack", func(m *nats.Msg) {
		natscommon.LogMessageReceived(m)
		log.Printf("Received ack from %s.  Total acks: %d", m.Reply, tracker.K8sPiCount)
		tracker.K8sPiCount++
		if tracker.K8sPiCount == totalK8sPis {
			log.Println("All k8s pis have acked.  Shutting down Pi Switch")
			ShutdownPiSwitch()
		}
	})

	log.Println("Subscription setup complete.  Polling battery status.")
	for {
		batteryStatus := pollBatteryStatus();
		log.Printf("Battery status: %#v", batteryStatus)
		if batteryStatus.IsOnBattery && batteryStatus.Percent <= 95 {
			if !tracker.Group1Deactivated {
				log.Println("Deactivating group 1")
				nc.Publish("ups.office", []byte("group1"))
				tracker.Group1Deactivated = true
			}
			if !tracker.Group2Deactivated && batteryStatus.Percent <=85 {
				log.Println("Deactivating group 2")
				nc.Publish("ups.office", []byte("group2"))
				tracker.Group2Deactivated = true
			}
			if !tracker.Group3Deactivated && batteryStatus.Percent <=40 {
				log.Println("Deactivating group 3")
				ShutdownPiSwitch()
				ShutdownNas()
				tracker.Group3Deactivated = true
			}
			if !tracker.Group4Deactivated && batteryStatus.Percent <=20 {
				log.Println("Deactivating group 4")
				ShutdownHass()
				tracker.Group4Deactivated = true
			}
		} else {
			// Exit
			log.Println("Battery is charged.  Turning things back on.")
			BringUpHass()
			if !tracker.NasIsOff {
				log.Println("NAS is not off, turning its dependents back on")
				BringUpPiSwitch()
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}


func pollBatteryStatus() BatteryStatus {
	return BatteryStatus{
		IsOnBattery: false,
		Percent: 100,
	}
}