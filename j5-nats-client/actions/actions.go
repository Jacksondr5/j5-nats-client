package actions

import (
	"log"
	"os/exec"
	"strings"
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
	cmd := exec.Command("gitlab-ctl", "graceful-kill")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Error shutting down GitLab")
		log.Println(err)
		log.Printf("Shutting down system anyway")
	} else {
		log.Println("GitLab shutdown complete successfully")
	}
	log.Printf("GitLab shutdown stdout: %s", out.String())

	Shutdown(nc, hostname)
}