package actions

import (
	"os/exec"
	"strings"
	"syscall"

	"github.com/kardianos/service"
	"github.com/nats-io/nats.go"
)

func Pong(nc *nats.Conn, logger service.Logger, hostname string) {
	logger.Info("Responding to ping")
	nc.Publish("pong", []byte(hostname))
}

func Shutdown(nc *nats.Conn, logger service.Logger, hostname string) {
	logger.Info("Shutting down system")
	nc.Publish("shutdown", []byte(hostname))
	nc.Drain()
	syscall.Exec("/sbin/shutdown", []string{"shutdown", "now"}, []string{})
}

func ShutdownGitLab(nc *nats.Conn, logger service.Logger, hostname string) {
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

	Shutdown(nc, logger, hostname)
}