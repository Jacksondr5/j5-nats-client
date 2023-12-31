package logic_test

import (
	"errors"
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/stretchr/testify/mock"
)

type mockKillableDevice struct {
	mock.Mock
}

func (m *mockKillableDevice) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockKillableDevice) IsOff() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockKillableDevice) SetIsOff(isOff bool) {
	m.Called(isOff)
}

func (m *mockKillableDevice) TurnOff() (err error) {
	args := m.Called()
	return args.Error(0)
}

func TestTurnOffDevice_HappyPath(t *testing.T) {
	// Given
	mockKillableDevice := &mockKillableDevice{}
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
	mockKillableDevice := &mockKillableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(true)
	
	// When
	logic.TurnOffDevice(mockKillableDevice)
	
	// Then
	mockKillableDevice.AssertExpectations(t)
}

func TestTurnOffDevice_Error(t *testing.T) {
	// Given
	mockKillableDevice := &mockKillableDevice{}
	mockKillableDevice.On("Name").Return("name")
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("TurnOff").Return(errors.New("error"))
	
	// When
	logic.TurnOffDevice(mockKillableDevice)
	
	// Then
	mockKillableDevice.AssertExpectations(t)
}