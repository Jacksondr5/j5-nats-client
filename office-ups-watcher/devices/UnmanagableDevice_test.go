package devices_test

import (
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/stretchr/testify/assert"
)

func TestNewUnmanageableDevice(t *testing.T) {
	// Given
	ud := devices.NewUnmanageableDevice("test")
	
	// When
	name := ud.Name()
	isOff := ud.IsOff()
	
	// Then
	assert.Equal(t, "test", name)
	assert.False(t, isOff)
}

func TestUnmanageableDevice_IsOff(t *testing.T) {
	// Given
	ud := devices.NewUnmanageableDevice("test")
	
	// When
	ud.SetIsOff(true)
	
	// Then
	assert.True(t, ud.IsOff())
}