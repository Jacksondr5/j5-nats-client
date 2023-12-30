package main

import (
	"log"
	"time"

	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/nats-io/nats.go"
)


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
	// Open a file in ./ups-watcher.log directory with the current time as the name
	// This will create a new file every time the program is run
	// This is useful for debugging
	// logFileName := fmt.Sprintf("./ups-watcher.log/%s", time.Now().Format("2006-01-02T15:04:05"))
	// logFile, err := os.Create(logFileName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)
	
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
	badBatteryStatusCount := 0
	isActive := false
	sleepTime := 10 * time.Second
	for {
		batteryStatus, err := battery.PollBatteryStatus();
		if err != nil {
			log.Println(err)
			log.Println("Error polling battery status, trying again....")
			badBatteryStatusCount++
			if badBatteryStatusCount > 5 {
				log.Println("Too many errors polling battery status, exiting.")
				break
			}
			time.Sleep(sleepTime)
			continue
		}
		log.Printf("Battery status: %#v", batteryStatus)
		if batteryStatus.IsOnBattery && batteryStatus.Percent <= 95 {
			log.Println("System is on battery.  Turning things off.")
			sleepTime = 1 * time.Second
			isActive = true
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
		} else if isActive && !batteryStatus.IsOnBattery && batteryStatus.Percent >= 95 {
			// Exit
			log.Println("Battery is charged.  Turning things back on.")
			if !tracker.NasIsOff {
				log.Println("NAS is not off, turning its dependents back on")
				BringUpPiSwitch()
			} else {
				log.Println("NAS is off, not turning its dependents back on")
			}
			break
		} else {
			log.Println("Everything is normal")
		}
		time.Sleep(sleepTime)
	}
}


