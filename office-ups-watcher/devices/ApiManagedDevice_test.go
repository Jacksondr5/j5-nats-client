package devices_test

import (
	"testing"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/devices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewApiManagedDevice(t *testing.T) {
	// Given
	ud := devices.NewApiManagedDevice("name", "", "", "", "", nil)
	
	// When
	name := ud.Name()
	isOff := ud.IsOff()
	
	// Then
	assert.Equal(t, "name", name)
	assert.False(t, isOff)
}

func TestApiManagedDevice_IsOff(t *testing.T) {
	// Given
	ud := devices.NewApiManagedDevice("", "", "", "", "", nil)
	
	// When
	ud.SetIsOff(true)
	
	// Then
	assert.True(t, ud.IsOff())
}

type mockHttpPost struct {
	mock.Mock
}

func (m *mockHttpPost) Post(url string, body string, name string, accessToken string) (err error) {
	args := m.Called(url, body, name, accessToken)
	return args.Error(0)
}

func TestApiManagedDevice_TurnOff(t *testing.T) {
	// Given
	mockHttpPost := &mockHttpPost{}
	ud := devices.NewApiManagedDevice("name", "access token", "turn off url", "", "body", mockHttpPost)
	mockHttpPost.On("Post", "turn off url", "body", "name shutdown", "access token").Return(nil)
	
	// When
	err := ud.TurnOff()
	
	// Then
	assert.Nil(t, err)
	mockHttpPost.AssertExpectations(t)
}

func TestApiManagedDevice_TurnOn(t *testing.T) {
	// Given
	mockHttpPost := &mockHttpPost{}
	ud := devices.NewApiManagedDevice("name", "access token", "", "turn on url", "body", mockHttpPost)
	mockHttpPost.On("Post", "turn on url", "body", "name startup", "access token").Return(nil)
	
	// When
	err := ud.TurnOn()
	
	// Then
	assert.Nil(t, err)
	mockHttpPost.AssertExpectations(t)
}