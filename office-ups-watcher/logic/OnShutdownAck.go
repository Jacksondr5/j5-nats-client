package logic

import (
	"time"

	"github.com/jacksondr5/go-monorepo/logger"
	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)



func OnShutdownAck(
	m *nats.Msg, 
	ackCounterChannel chan<- PiMonitorAction,
	monitorChannel <-chan int,
) {
	natscommon.LogMessageReceived(m)
	logger.InfoWithFields("Received ack", log.Fields{
		"ack": string(m.Data),
	})
	ackCounterChannel <- PiMonitorAction{
		Hostname: string(m.Data),
		EverythingShutDown: false,
		Timestamp: time.Now().Unix(),
	}
	<-monitorChannel
}