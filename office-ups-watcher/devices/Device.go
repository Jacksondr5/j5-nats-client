package devices

type Device interface {
	Name() string
	IsOff() bool
	SetIsOff(isOff bool)
}