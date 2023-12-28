package main

import (
	"log"
	"os"
	"time"

	"github.com/jacksondr5/j5-nats-client/actions"
	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

type Subscription struct {
	Action string
	Name string
	Subject string
}

type Config struct {
	Name string
	Subscriptions []Subscription
}

var hostname string

func main() {
	config := ReadYamlFile("j5-nats-client-config.yaml")
	log.Printf("%#v",config)

	hostname = getHostname()
	log.Println("Starting up")
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222")
	defer nc.Drain()
	log.Println("Connected to NATS")

	log.Println("Setting up subscriptions")
	for i := range config.Subscriptions {
		subscription := config.Subscriptions[i]
		log.Printf("Setting up subscription for \"%s\" on topic \"%s\" with action \"%s\"", subscription.Name, subscription.Subject, subscription.Action)
		nc.Subscribe(subscription.Subject, func(m *nats.Msg) {
			logMessageReceived(m)
			log.Printf("Executing action \"%s\" in response to message on subject \"%s\"", subscription.Action, m.Subject)
			switch subscription.Action {
			case "pong":
				actions.Pong(nc, hostname)
			case "shutdown":
				actions.Shutdown(nc, hostname)
			case "shutdown-gitlab":
				actions.ShutdownGitLab(nc, hostname)
			default:
				log.Printf("Unknown action %s", subscription.Action)
			}
		})
	}


	log.Println("Subscription setup complete.  Waiting for events.")

	for {
		time.Sleep(1 * time.Second)
	}
}

func logMessageReceived(m *nats.Msg) {
	log.Printf("Received a message from subject %s: %s \n", m.Subject, string(m.Data))
}

func getHostname() string {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		panic(hostnameErr)
	}
	return hostname
}

func ReadYamlFile(path string) *Config {
	buf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	return &config
}