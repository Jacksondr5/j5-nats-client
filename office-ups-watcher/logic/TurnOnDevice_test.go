package logic_test

import (
	"errors"
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/stretchr/testify/mock"
)

type mockRaisableDevice struct {
	mock.Mock
}

func (m *mockRaisableDevice) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRaisableDevice) IsOff() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockRaisableDevice) SetIsOff(isOff bool) {
	m.Called(isOff)
}

func (m *mockRaisableDevice) TurnOn() (err error) {
	args := m.Called()
	return args.Error(0)
}

func TestTurnOnDevice_HappyPath(t *testing.T) {
	// Given
	mockRaisableDevice := &mockRaisableDevice{}
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
	mockRaisableDevice := &mockRaisableDevice{}
	mockRaisableDevice.On("Name").Return("name")
	mockRaisableDevice.On("IsOff").Return(false)
	
	// When
	logic.TurnOnDevice(mockRaisableDevice)
	
	// Then
	mockRaisableDevice.AssertExpectations(t)
}

func TestTurnOnDevice_TurnOnError(t *testing.T) {
	// Given
	mockRaisableDevice := &mockRaisableDevice{}
	mockRaisableDevice.On("Name").Return("name")
	mockRaisableDevice.On("IsOff").Return(true)
	mockRaisableDevice.On("TurnOn").Return(errors.New("error"))
	
	// When
	logic.TurnOnDevice(mockRaisableDevice)
	
	// Then
	mockRaisableDevice.AssertExpectations(t)
}