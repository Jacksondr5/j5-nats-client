package logic_test

import (
	"errors"
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	logicTest "github.com/jacksondr5/go-monorepo/office-ups-watcher/logic/test"
)

func TestTurnOffDevice_HappyPath(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("SetIsOff", true).Return()
	mockKillableDevice.On("TurnOff").Return(nil)
	
	// When
	logic.TurnOffDevice(mockKillableDevice)
	
	// Then
	mockKillableDevice.AssertExpectations(t)
}

func TestTurnOffDevice_AlreadyOff(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(true)
	
	// When
	logic.TurnOffDevice(mockKillableDevice)
	
	// Then
	mockKillableDevice.AssertExpectations(t)
}

func TestTurnOffDevice_Error(t *testing.T) {
	// Given
	mockKillableDevice := &logicTest.MockManageableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("TurnOff").Return(errors.New("error"))
	
	// When
	logic.TurnOffDevice(mockKillableDevice)
	
	// Then
	mockKillableDevice.AssertExpectations(t)
}