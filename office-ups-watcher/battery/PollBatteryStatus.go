package battery

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jacksondr5/go-monorepo/office-ups-watcher/logger"
)

type BatteryStatus struct {
	IsOnBattery bool
	Percent int
}

type BatteryPollerImpl struct {}

// Toggle this for testing
const upscCmd = "upsc"
// const upscCmd = "./upsc.sh"

func (bp BatteryPollerImpl) PollBatteryStatus() (BatteryStatus, error) {
	charge, chargeErr := execCmd(upscCmd, "cyberpower@localhost", "battery.charge")
	status, statusErr := execCmd(upscCmd, "cyberpower@localhost", "ups.status")
	if chargeErr != nil || statusErr != nil {
		if chargeErr != nil {
			logger.Error("error getting battery charge", chargeErr)
		}
		if statusErr != nil {
			logger.Error("error getting battery status", statusErr)
		}
		return BatteryStatus{}, errors.New("error polling battery status")
	}

	var isOnBattery bool
	if strings.Contains(status, "OB") {
		isOnBattery = true
	} else if strings.Contains(status, "OL") {
		isOnBattery = false
	} else {
		return BatteryStatus{}, errors.New("battery status not recognized")
	}

	charge = strings.TrimSuffix(charge, "\n")
	percent, err := strconv.Atoi(charge)
	if err != nil {
		return BatteryStatus{}, errors.New("error converting charge to int")
	}

	return BatteryStatus{
		IsOnBattery: isOnBattery,
		Percent: percent,
	}, nil
}

func execCmd(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}