package logic_test

import (
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestOnShutdownAck(t *testing.T) {
	// Given
	mockMessage := nats.NewMsg("subject")
	mockMessage.Data = []byte("name")
	mockAckChannel := make(chan logic.PiMonitorAction)
	mockMonitorChannel := make(chan int)
	expected := logic.PiMonitorAction{
		EverythingShutDown: false,
		Hostname: "name",
		Timestamp: 0,
	}
	go func() {
		test := <-mockAckChannel
		mockMonitorChannel <- 1
		// Then
		assert.Equal(t, expected.EverythingShutDown, test.EverythingShutDown)
		assert.Equal(t, expected.Hostname, test.Hostname)
		assert.Greater(t, test.Timestamp, expected.Timestamp)
	}()

	// When
	logic.OnShutdownAck(mockMessage, mockAckChannel, mockMonitorChannel)
}
