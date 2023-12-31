package logic_test

import (
	"errors"
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	logicTest "github.com/jacksondr5/go-monorepo/office-ups-watcher/logic/test"
)


func TestTurnOnDevice_HappyPath(t *testing.T) {
	// Given
	mockRaisableDevice := &logicTest.MockManageableDevice{}
	mockRaisableDevice.On("Name").Return("name")
	mockRaisableDevice.On("IsOff").Return(true)
	mockRaisableDevice.On("SetIsOff", false).Return()
	mockRaisableDevice.On("TurnOn").Return(nil)
	
	// When
	logic.TurnOnDevice(mockRaisableDevice)
	
	// Then
	mockRaisableDevice.AssertExpectations(t)
}

func TestTurnOnDevice_AlreadyOn(t *testing.T) {
	// Given
	mockRaisableDevice := &logicTest.MockManageableDevice{}
	mockRaisableDevice.On("Name").Return("name")
	mockRaisableDevice.On("IsOff").Return(false)
	
	// When
	logic.TurnOnDevice(mockRaisableDevice)
	
	// Then
	mockRaisableDevice.AssertExpectations(t)
}

func TestTurnOnDevice_TurnOnError(t *testing.T) {
	// Given
	mockRaisableDevice := &logicTest.MockManageableDevice{}
	mockRaisableDevice.On("Name").Return("name")
	mockRaisableDevice.On("IsOff").Return(true)
	mockRaisableDevice.On("TurnOn").Return(errors.New("error"))
	
	// When
	logic.TurnOnDevice(mockRaisableDevice)
	
	// Then
	mockRaisableDevice.AssertExpectations(t)
}