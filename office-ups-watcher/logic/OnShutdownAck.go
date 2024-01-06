package logic

import (
	"fmt"

	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logger"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

const totalK8sPis = 11

func OnShutdownAck(m *nats.Msg, devices *ManagedDevices, k8sPiCount int) int {
	natscommon.LogMessageReceived(m)
	logger.InfoWithFields("Received ack", log.Fields{
		"ack": string(m.Data),
		"totalAcks": k8sPiCount,
	})
	k8sPiCount++
	if k8sPiCount == totalK8sPis {
		logger.Info("All k8s pis have acked.  Shutting down Pi Switch")
		TurnOffDevice(devices.PiSwitch)
	} else if k8sPiCount > totalK8sPis {
		logger.Warning(fmt.Sprintf("More than %d k8s pis have acked.  THIS IS NOT EXPECTED!  Shutting down Pi Switch anyway...", totalK8sPis))
		TurnOffDevice(devices.PiSwitch)
	}
	return k8sPiCount
}