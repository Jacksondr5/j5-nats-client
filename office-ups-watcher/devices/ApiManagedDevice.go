package devices

import (
	"fmt"
)

type HttpClient interface {
	Post(url string, bodyData string, requestName string, accessToken string) error
}

type ApiManagedDevice struct {
	accessToken string
	isOff bool
	httpClient HttpClient
	name string
	requestBody string
	turnOffUrl string
	turnOnUrl string
}

func (amd ApiManagedDevice) IsOff() bool {
	return amd.isOff
}

func (amd ApiManagedDevice) Name() string {
	return amd.name
}

func (amd *ApiManagedDevice) SetIsOff(isOff bool) {
	amd.isOff = isOff
}

func (amd *ApiManagedDevice) TurnOff() (err error) {
	return amd.httpClient.Post(
		amd.turnOffUrl,
		amd.requestBody,
		fmt.Sprintf("%s shutdown", amd.name),
		amd.accessToken,
	)
}

func (amd *ApiManagedDevice) TurnOn() (err error) {
	return amd.httpClient.Post(
		amd.turnOnUrl,
		amd.requestBody,
		fmt.Sprintf("%s startup", amd.name),
		amd.accessToken,
	)
}

func NewApiManagedDevice(
	name string, 
	accessToken string, 
	turnOffUrl string, 
	turnOnUrl string, 
	requestBody string,
	httpClient HttpClient,
) *ApiManagedDevice {
	return &ApiManagedDevice{
		accessToken: accessToken,
		isOff: false,
		httpClient: httpClient,
		name: name,
		requestBody: requestBody,
		turnOffUrl: turnOffUrl,
		turnOnUrl: turnOnUrl,
	}
}