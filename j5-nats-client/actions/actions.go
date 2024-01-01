package actions

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kardianos/service"
	"github.com/nats-io/nats.go"
)

func Pong(nc *nats.Conn, logger service.Logger, hostname string) {
	logger.Info("Responding to ping")
	nc.Publish("pong", []byte(hostname))
}

func ShutdownUbuntu(nc *nats.Conn, logger service.Logger, hostname string, subject string) {
	logger.Info("Shutting down Ubuntu system")
	shutdownSystem(
		nc,
		logger,
		hostname,
		subject,
		"/sbin/shutdown", 
		[]string{"shutdown", "now"},
	)
}

func ShutdownWindows(nc *nats.Conn, logger service.Logger, hostname string, subject string) {
	logger.Info("Shutting down Windows system")
	shutdownSystem(
		nc,
		logger,
		hostname,
		subject,
		"shutdown", 
		[]string{
			"/s", 
			"/t", "0", 
			"/c", "\"Shutting down due to J5 NATS Client\"", 
			"/d", "u:6:12",
		},
	)
	// user32 := syscall.MustLoadDLL("user32")
	// defer user32.Release()

	// exitwin := user32.MustFindProc("ExitWindowsEx")

	// r1, _, err := exitwin.Call(0x08, 0)
	// if r1 != 1 {
	// 	fmt.Println("Failed to initiate shutdown:", err)
	// }
}

func shutdownSystem(
	nc *nats.Conn, 
	logger service.Logger, 
	hostname string, 
	subject string, 
	command string, 
	args []string,
) {
	publishShutdownAck(nc, hostname, subject)
	nc.Drain()
	// err := syscall.Exec(command, args, []string{})
	cmd := exec.Command(command, args...)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Info("Error shutting down system")
		logger.Info(err)
	}
}

func ShutdownGitLabProcess(nc *nats.Conn, logger service.Logger, hostname string, subject string) {
	logger.Info("Shutting down GitLab")
	cmd := exec.Command("gitlab-ctl", "graceful-kill")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Info("Error shutting down GitLab")
		logger.Info(err)
		logger.Info("Shutting down system anyway")
	} else {
		logger.Info("GitLab shutdown complete successfully")
	}
	logger.Infof("GitLab shutdown stdout: %s", out.String())

	ShutdownUbuntu(nc, logger, hostname, subject)
}

func publishShutdownAck(nc *nats.Conn, hostname string, subject string) {
	nc.Publish(fmt.Sprintf("%s.ack", subject), []byte(hostname))
}