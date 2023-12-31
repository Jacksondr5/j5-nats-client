package test

import "github.com/stretchr/testify/mock"

type MockManageableDevice struct {
	mock.Mock
}

func (m *MockManageableDevice) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockManageableDevice) IsOff() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockManageableDevice) SetIsOff(isOff bool) {
	m.Called(isOff)
}

func (m *MockManageableDevice) TurnOff() (err error) {
	args := m.Called()
	return args.Error(0)
}

func (m *MockManageableDevice) TurnOn() (err error) {
	args := m.Called()
	return args.Error(0)
}