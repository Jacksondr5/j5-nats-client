package logic

import (
	"log"

	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/nats-io/nats.go"
)

const totalK8sPis = 11

func OnShutdownAck(m *nats.Msg, devices *ManagedDevices, k8sPiCount int) int {
	natscommon.LogMessageReceived(m)
	log.Printf("Received ack from %s.  Total acks: %d", string(m.Data), k8sPiCount)
	k8sPiCount++
	if k8sPiCount == totalK8sPis {
		log.Println("All k8s pis have acked.  Shutting down Pi Switch")
		TurnOffDevice(devices.PiSwitch)
	} else if k8sPiCount > totalK8sPis {
		log.Printf("More than %d k8s pis have acked.  THIS IS NOT EXPECTED!  Shutting down Pi Switch anyway...", totalK8sPis)
		TurnOffDevice(devices.PiSwitch)
	}
	return k8sPiCount
}