package natscommon

import (
	"log"

	"github.com/nats-io/nats.go"
)

func LogMessageReceived(m *nats.Msg) {
	log.Printf("Received a message from subject %s: %s \n", m.Subject, string(m.Data))
}