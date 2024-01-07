package logic_test

import (
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logic"
	logicTest "github.com/jacksondr5/go-monorepo/office-ups-watcher/logic/test"
)

func TestMonitorPiShutdown_SingleAck(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := &logicTest.MockManageableDevice{}
	message :=logic.PiMonitorAction{
		EverythingShutDown: false, 
		Hostname: "name", 
		Timestamp: 0,
	}

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, 2)
	actionFeed <- message
	<-monitor

	// Then
	// Shouldn't do anything since we only have one ack and there are 2 total Pis
	mockKillableDevice.AssertExpectations(t)
}

func TestMonitorPiShutdown_AllPisAcked(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := getMockDevice("name", false, true, "TurnOff")
	totalPis := 2

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, totalPis)
	for i := 0; i < totalPis; i++ {
		actionFeed <- logic.PiMonitorAction{
			EverythingShutDown: false,
			Hostname: string(rune(i)),
			Timestamp: 0,
		}
		<-monitor
	}

	// Then
	// Expect the device to be turned off since all pis have acked
	mockKillableDevice.AssertExpectations(t)
}

func TestMonitorPiShutdown_MoreThanAllPisAck(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := &logicTest.MockManageableDevice{}
	// Because we send one more than the total number of pis, we expect the functions to be called twice
	mockKillableDevice.On("Name").Return("name").Twice()
	mockKillableDevice.On("IsOff").Return(false)
	mockKillableDevice.On("IsOff").Return(true)
	mockKillableDevice.On("SetIsOff", true).Return()
	mockKillableDevice.On("TurnOff").Return(nil)
	totalPis := 2

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, totalPis)
	for i := 0; i < totalPis + 1; i++ {
		actionFeed <- logic.PiMonitorAction{
			EverythingShutDown: false,
			Hostname: string(rune(i)),
			Timestamp: 0,
		}
		<-monitor
	}

	// Then
	mockKillableDevice.AssertExpectations(t)
}

func TestMonitorPiShutdown_EverythingShutDown(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := &logicTest.MockManageableDevice{}
	message := logic.PiMonitorAction{
		EverythingShutDown: true,
		Hostname: "",
		Timestamp: 0,
	}

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, 1)
	actionFeed <- message
	<-monitor

	// Then
	// Nothing should be turned off here
	mockKillableDevice.AssertExpectations(t)
}

func TestMonitorPiShutdown_OldDataInChannel(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := &logicTest.MockManageableDevice{}
	everythingShutDownMessage := logic.PiMonitorAction{
		EverythingShutDown: true,
		Hostname: "",
		Timestamp: 1,
	}
	oldMessage := logic.PiMonitorAction{
		EverythingShutDown: false,
		Hostname: "this message was generated before the shutdown message",
		Timestamp: 0,
	}

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, 1)
	actionFeed <- everythingShutDownMessage
	<-monitor
	actionFeed <- oldMessage
	<-monitor

	// Then
	// Expect nothing to be shut down since the ack message was generated before the everything shutdown message
	mockKillableDevice.AssertExpectations(t)
}

func TestMonitorPiShutdown_EverythingShutDownAndThenNewAck(t *testing.T) {
	// Given
	actionFeed := make(chan logic.PiMonitorAction)
	monitor := make(chan int)
	mockKillableDevice := &logicTest.MockManageableDevice{}
	// Because we reset the counter, we expect the device to get shut down twice
	mockKillableDevice.On("Name").Return("name").Twice()
	mockKillableDevice.On("IsOff").Return(false).Twice()
	mockKillableDevice.On("SetIsOff", true).Return().Twice()
	mockKillableDevice.On("TurnOff").Return(nil).Twice()
	oldAckMessage := logic.PiMonitorAction{
		EverythingShutDown: false,
		Hostname: "this message was generated before the shutdown message",
		Timestamp: 0,
	}
	everythingShutDownMessage := logic.PiMonitorAction{
		EverythingShutDown: true,
		Hostname: "",
		Timestamp: 1,
	}
	newAckMessage := logic.PiMonitorAction{
		EverythingShutDown: false,
		Hostname: "this message was generated after the shutdown message",
		Timestamp: 2,
	}

	// When
	go logic.MonitorPiShutdown(actionFeed, monitor, &logic.ManagedDevices{PiSwitch: mockKillableDevice}, 1)
	actionFeed <- oldAckMessage
	<-monitor
	actionFeed <- everythingShutDownMessage
	<-monitor
	actionFeed <- newAckMessage
	<-monitor

	// Then
	// Expect the device to be turned off since the ack message was generated after the everything shutdown message
	mockKillableDevice.AssertExpectations(t)
}
