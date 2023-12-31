package logic_test

import (
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	logicTest "github.com/jacksondr5/go-monorepo/office-ups-watcher/logic/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)


func TestOnShutdownAck_FirstAck(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockMessage := nats.NewMsg("subject")
	mockMessage.Data = []byte("name")
	k8sPisCount := 1

	// When
	newPiCount := logic.OnShutdownAck(mockMessage, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, k8sPisCount)

	// Then
	assert.Equal(t, 2, newPiCount)
	mockKillableDevice.AssertExpectations(t)
}

func TestOnShutdownAck_AllPisAcked(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("SetIsOff", true).Return()
	mockKillableDevice.On("TurnOff").Return(nil)
	mockMessage := nats.NewMsg("subject")
	mockMessage.Data = []byte("name")
	k8sPisCount := 11

	// When
	newPiCount := logic.OnShutdownAck(mockMessage, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, k8sPisCount)

	// Then
	assert.Equal(t, 12, newPiCount)
	mockKillableDevice.AssertExpectations(t)
}

func TestOnShutdownAck_MorePisThanExpectedAcked(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("SetIsOff", true).Return()
	mockKillableDevice.On("TurnOff").Return(nil)
	mockMessage := nats.NewMsg("subject")
	mockMessage.Data = []byte("name")
	k8sPisCount := 13

	// When
	newPiCount := logic.OnShutdownAck(mockMessage, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, k8sPisCount)

	// Then
	assert.Equal(t, 14, newPiCount)
	mockKillableDevice.AssertExpectations(t)
}