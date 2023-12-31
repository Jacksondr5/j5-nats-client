package devices

type UnmanageableDevice struct {
	isOff bool
	name string
}

func (ud UnmanageableDevice) IsOff() bool {
	return ud.isOff
}

func (ud UnmanageableDevice) Name() string {
	return ud.name
}

func (ud *UnmanageableDevice) SetIsOff(isOff bool) {
	ud.isOff = isOff
}

func NewUnmanageableDevice(name string) *UnmanageableDevice {
	return &UnmanageableDevice{
		isOff: false,
		name: name,
	}
}