package logic

import (
	"fmt"

	"github.com/jacksondr5/go-monorepo/logger"
	log "github.com/sirupsen/logrus"
)

type PiMonitorAction struct {
	Hostname string
	EverythingShutDown bool
	Timestamp int64
}

func MonitorPiShutdown(
	actionFeed <-chan PiMonitorAction, 
	ack chan<- int, 
	devices *ManagedDevices,
	totalK8sPis int,
) {
	k8sPiCount := 0
	lastShutdownTime := int64(0)
	for {
		action := <-actionFeed
		if action.Timestamp < lastShutdownTime {
			logger.InfoWithFields("Received old action to MonitorPiShutdown", log.Fields{
				"action": action,
				"lastShutdownTime": lastShutdownTime,
			})
			ack <- k8sPiCount
			continue
		} else if action.EverythingShutDown {
			logger.Info("Pi Switch has been turned off.  Resetting counter")
			k8sPiCount = 0
			lastShutdownTime = action.Timestamp
		} else {
			k8sPiCount++
		}
		logger.InfoWithFields("Received new action to MonitorPiShutdown", log.Fields{
			"action": action,
			"totalAcks": k8sPiCount,
		})

		if k8sPiCount == totalK8sPis {
			logger.Info("All k8s pis have acked.  Shutting down Pi Switch")
			TurnOffDevice(devices.PiSwitch)
		} else if k8sPiCount > totalK8sPis {
			logger.WarningWithFields(
				fmt.Sprintf(
					"More than %d k8s pis have acked.  THIS IS NOT EXPECTED!  Shutting down Pi Switch anyway...", 
					totalK8sPis,
				), 
				log.Fields{
					"ackCount": k8sPiCount,
					"totalPis": totalK8sPis,
				},
			)
			TurnOffDevice(devices.PiSwitch)
		}
		ack <- k8sPiCount
	}
}