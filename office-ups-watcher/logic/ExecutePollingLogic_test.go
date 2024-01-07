package logic_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/battery"
	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBatteryPoller struct {
	mock.Mock
}

func (m *mockBatteryPoller) PollBatteryStatus() (battery.BatteryStatus, error) {
	args := m.Called()
	return args.Get(0).(battery.BatteryStatus), args.Error(1)
}

type mockNatsConn struct {
	mock.Mock
	nats.Conn
}

func (m *mockNatsConn) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

type mockDevice struct {
	mock.Mock
}

func (m *mockDevice) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockDevice) IsOff() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockDevice) SetIsOff(isOff bool) {
	m.Called(isOff)
}

func (m *mockDevice) TurnOn() (err error) {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDevice) TurnOff() (err error) {
	args := m.Called()
	return args.Error(0)
}

func getMockBatteryPoller(status bool, percentage int) *mockBatteryPoller {
	mockBattery := &mockBatteryPoller{}
	mockBattery.On("PollBatteryStatus").Return(battery.BatteryStatus{IsOnBattery: status, Percent: percentage}, nil)
	return mockBattery
}

func getMockNatConn(group string) *mockNatsConn {
	mockNatsConn := &mockNatsConn{}
	mockNatsConn.On("Publish", "ups.office", []byte(group)).Return(nil)
	return mockNatsConn
}

func getMockDevice(name string, isOff bool, setIsOff bool, turnMethod string) *mockDevice {
	mockDevice := &mockDevice{}
	mockDevice.On("Name").Return(name)
	mockDevice.On("IsOff").Return(isOff)
	mockDevice.On("SetIsOff", setIsOff).Return()
	mockDevice.On(turnMethod).Return(nil)
	return mockDevice
}

func TestExecuteLogic_BatteryJustTurnedOn(t *testing.T) {
	// Given
	mockBattery := getMockBatteryPoller(true, 95)
	mockNatsConn := getMockNatConn("group1")
	mockDevice := &mockDevice{}
	mockDevices := &logic.ManagedDevices{PiSwitch: mockDevice, Nas: mockDevice}
	tracker := &logic.Tracker{IsActive: false}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When
	sleepTime := logic.ExecutePollingLogic(tracker, mockNatsConn, mockDevices, mockBattery, mockAckChannel)
	
	// Then
	assert.Equal(t, 1*time.Second, sleepTime)
	assert.True(t, tracker.IsActive)
	mockBattery.AssertExpectations(t)
	mockNatsConn.AssertExpectations(t)
	mockDevice.AssertExpectations(t)
}

func TestExecuteLogic_LinePowerCameBackOnAfterAShutdown(t *testing.T) {
	// Given
	mockBattery := getMockBatteryPoller(false, 100)
	mockNatsConn := &mockNatsConn{}
	mockDevice := &mockDevice{}
	mockDevice.On("IsOff").Return(false)
	mockDevice.On("Name").Return("name")
	mockDevices := &logic.ManagedDevices{PiSwitch: mockDevice, Nas: mockDevice}
	tracker := &logic.Tracker{
		IsActive: true, 
		Group1IsDeactivated: true, 
		Group2IsDeactivated: true, 
		Group3IsDeactivated: true,
	}
	mockAckChannel := make(chan logic.PiMonitorAction)

	// When
	sleepTime := logic.ExecutePollingLogic(tracker, mockNatsConn, mockDevices, mockBattery, mockAckChannel)

	// Then
	assert.Equal(t, 10*time.Second, sleepTime)
	assert.False(t, tracker.IsActive)
	assert.False(t, tracker.Group1IsDeactivated)
	assert.False(t, tracker.Group2IsDeactivated)
	assert.False(t, tracker.Group3IsDeactivated)
	mockBattery.AssertExpectations(t)
	mockNatsConn.AssertExpectations(t)
	mockDevice.AssertExpectations(t)
}

func TestExecuteLogic_LinePowerIsOnButBatteryNotCharged(t *testing.T) {
	// Given
	tracker := &logic.Tracker{IsActive: true}
	mockAckChannel := make(chan logic.PiMonitorAction)

	// When
	sleepTime := logic.ExecutePollingLogic(tracker, nil, nil, getMockBatteryPoller(false, 70), mockAckChannel)

	// Then
	assert.Equal(t, 10*time.Second, sleepTime)
	assert.True(t, tracker.IsActive)
}

func TestExecuteLogic_LinePowerIsOnAndBatteryChargedAndNasIsOn(t *testing.T) {
	// Given
	mockNas := &mockDevice{}
	mockNas.On("IsOff").Return(false)
	mockPiSwitch :=  getMockDevice("pi switch", true, false, "TurnOn")
	mockDevices := &logic.ManagedDevices{PiSwitch: mockPiSwitch, Nas: mockNas}
	tracker := &logic.Tracker{IsActive: true}
	mockAckChannel := make(chan logic.PiMonitorAction)

	// When
	sleepTime := logic.ExecutePollingLogic(tracker, nil, mockDevices, getMockBatteryPoller(false, 100), mockAckChannel)

	// Then
	assert.Equal(t, 10*time.Second, sleepTime)
	assert.False(t, tracker.IsActive)
	mockPiSwitch.AssertExpectations(t)
}

func TestExecuteLogic_LinePowerIsOnAndBatteryChargedAndNasIsOff(t *testing.T) {
	// Given
	mockNas := &mockDevice{}
	mockNas.On("IsOff").Return(true)
	mockPiSwitch := &mockDevice{}
	mockDevices := &logic.ManagedDevices{PiSwitch: mockPiSwitch, Nas: mockNas}
	tracker := &logic.Tracker{IsActive: true}
	mockAckChannel := make(chan logic.PiMonitorAction)

	// When
	sleepTime := logic.ExecutePollingLogic(tracker, nil, mockDevices, getMockBatteryPoller(false, 100), mockAckChannel)

	// Then
	assert.Equal(t, 10*time.Second, sleepTime)
	assert.False(t, tracker.IsActive)
	// Ensure nothing on the PiSwitch was called
	mockPiSwitch.AssertExpectations(t)
}

func TestExecuteLogic_EverythingIsNormal(t *testing.T) {
	// Given
	tracker := &logic.Tracker{IsActive: false}
	mockAckChannel := make(chan logic.PiMonitorAction)

	// When
	sleepTime := logic.ExecutePollingLogic(tracker, nil, nil, getMockBatteryPoller(false, 70), mockAckChannel)

	// Then
	assert.Equal(t, 10*time.Second, sleepTime)
	assert.False(t, tracker.IsActive)
}

func TestExecuteLogic_BatteryPercentageBelow95AndGroup1IsNotOff(t *testing.T) {
	// Given
	mockNatsConn := getMockNatConn("group1")
	tracker := &logic.Tracker{IsActive: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, mockNatsConn, nil, getMockBatteryPoller(true, 94), mockAckChannel)
	
	// Then
	assert.True(t, tracker.Group1IsDeactivated)
	mockNatsConn.AssertExpectations(t)
}

func TestExecuteLogic_BatteryPercentageBelow95AndGroup1IsAlreadyOff(t *testing.T) {
	// Given
	mockNatsConn := &mockNatsConn{}
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, mockNatsConn, nil, getMockBatteryPoller(true, 94), mockAckChannel)
	
	// Then
	// Expect nothing to be called on the nats connection
	mockNatsConn.AssertExpectations(t)
}

func TestExecuteLogic_BatteryPercentageBelow85AndGroup2IsNotOff(t *testing.T) {
	// Given
	mockNatsConn := getMockNatConn("group2")
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, mockNatsConn, nil, getMockBatteryPoller(true, 84), mockAckChannel)
	
	// Then
	assert.True(t, tracker.Group2IsDeactivated)
	mockNatsConn.AssertExpectations(t)
}

func TestExecuteLogic_BatteryPercentageBelow85AndGroup2IsAlreadyOff(t *testing.T) {
	// Given
	mockNatsConn := &mockNatsConn{}
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true, Group2IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, mockNatsConn, nil, getMockBatteryPoller(true, 84), mockAckChannel)
	
	// Then
	// Expect nothing to be called on the nats connection
	mockNatsConn.AssertExpectations(t)
}

func TestExecuteLogic_BatteryPercentageBelow40AndGroup3IsNotOff(t *testing.T) {
	// Given
	mockNatsConn := &mockNatsConn{}
	mockPiSwitch := getMockDevice("pi switch", false, true, "TurnOff")
	mockNas := getMockDevice("nas", false, true, "TurnOff")
	mockDevices := &logic.ManagedDevices{PiSwitch: mockPiSwitch, Nas: mockNas}
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true, Group2IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	expectedAction := logic.PiMonitorAction{
		EverythingShutDown: true,
		Hostname: "",
		Timestamp: 0,
	}
	
	// When, _
	go logic.ExecutePollingLogic(tracker, mockNatsConn, mockDevices, getMockBatteryPoller(true, 39), mockAckChannel)
	acks := <-mockAckChannel
	
	// Then
	assert.True(t, tracker.Group3IsDeactivated)
	assert.Equal(t, expectedAction.EverythingShutDown, acks.EverythingShutDown)
	assert.Equal(t, expectedAction.Hostname, acks.Hostname)
	assert.Greater(t, acks.Timestamp, expectedAction.Timestamp)
	mockPiSwitch.AssertExpectations(t)
	mockNas.AssertExpectations(t)
}

func TestExecuteLogic_BatteryPercentageBelow40AndGroup3IsAlreadyOff(t *testing.T) {
	// Given
	mockNatsConn := &mockNatsConn{}
	mockPiSwitch := &mockDevice{}
	mockNas := &mockDevice{}
	mockDevices := &logic.ManagedDevices{PiSwitch: mockPiSwitch, Nas: mockNas}
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true, Group2IsDeactivated: true, Group3IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, mockNatsConn, mockDevices, getMockBatteryPoller(true, 39), mockAckChannel)
	
	// Then
	mockPiSwitch.AssertExpectations(t)
	mockNas.AssertExpectations(t)
}

func TestExecuteLogic_GettingBatteryStatusFailed(t *testing.T) {
	// Given
	mockNatsConn := &mockNatsConn{}
	mockBattery := &mockBatteryPoller{}
	mockBattery.On("PollBatteryStatus").Return(battery.BatteryStatus{}, errors.New("error"))
	tracker := &logic.Tracker{IsActive: true, Group1IsDeactivated: true, Group2IsDeactivated: true, Group3IsDeactivated: true}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	// When, _
	logic.ExecutePollingLogic(tracker, nil, nil, mockBattery, mockAckChannel)
	
	// Then
	assert.Equal(t, 1, tracker.BadBatteryStatusCount)
	mockNatsConn.AssertExpectations(t)
	mockBattery.AssertExpectations(t)
}

func FuzzExecuteLogic(f *testing.F) {
	mockNatsConn := &mockNatsConn{}
	mockNatsConn.On("Publish", "ups.office", []byte("group1")).Return(nil)
	mockNatsConn.On("Publish", "ups.office", []byte("group2")).Return(nil)

	mockOnPiSwitch := getMockDevice("pi switch", false, true, "TurnOff")
	// mockOffPiSwitch := getMockDevice("pi switch", true, false, "TurnOn")
	mockOnNas := getMockDevice("nas", false, true, "TurnOff")
	// mockOffNas := getMockDevice("nas", true, false, "TurnOn")
	mockDevices := &logic.ManagedDevices{PiSwitch: mockOnPiSwitch, Nas: mockOnNas}
	mockAckChannel := make(chan logic.PiMonitorAction)
	
	f.Add(true, false, false, false, false, uint8(100))

	f.Fuzz(func(
		t *testing.T, 
		isActive bool, 
		group1IsDeactivated bool, 
		group2IsDeactivated bool, 
		group3IsDeactivated bool, 
		isOnBattery bool, 
		percentage uint8,
	) {
		mockBattery := getMockBatteryPoller(isOnBattery, int(percentage)) 
		logic.ExecutePollingLogic(
 			&logic.Tracker{
				IsActive: isActive,
				Group1IsDeactivated: group1IsDeactivated,
				Group2IsDeactivated: group2IsDeactivated,
				Group3IsDeactivated: group3IsDeactivated,
			}, 
			mockNatsConn, 
			mockDevices, 
			mockBattery,
			mockAckChannel,
		)
	})
}