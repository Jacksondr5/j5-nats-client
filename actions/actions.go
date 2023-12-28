package actions

import (
	"log"
	"syscall"

	"github.com/nats-io/nats.go"
)

func Pong(nc *nats.Conn) {
	log.Println("Responding to ping")
	nc.Publish("pong", []byte("pong"))
}

func Shutdown(nc *nats.Conn) {
	log.Println("Shutting down")
	// TODO: Figure out how to communicate shutdown to network
	// nc.Publish("shutdown", []byte("shutdown"))
	nc.Drain()
	// exec.Command("sudo", "shutdown now").Run()
	syscall.Exec("/sbin/shutdown", []string{"shutdown", "now"}, []string{})
}

func ShutdownGitLab(nc *nats.Conn) {
	log.Println("Shutting down GitLab")
}