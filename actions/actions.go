package actions

import (
	"log"
	"syscall"

	"github.com/nats-io/nats.go"
)

func Pong(nc *nats.Conn, hostname string) {
	log.Println("Responding to ping")
	nc.Publish("pong", []byte(hostname))
}

func Shutdown(nc *nats.Conn, hostname string) {
	log.Println("Shutting down system")
	nc.Publish("shutdown", []byte(hostname))
	nc.Drain()
	syscall.Exec("/sbin/shutdown", []string{"shutdown", "now"}, []string{})
}

func ShutdownGitLab(nc *nats.Conn, hostname string) {
	log.Println("Shutting down GitLab")
	syscall.Exec("/usr/bin/gitlab-ctl", []string{"gitlab-ctl", "stop"}, []string{})
	Shutdown(nc, hostname)
}