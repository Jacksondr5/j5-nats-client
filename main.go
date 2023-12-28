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
	Topic string
}

type Config struct {
	Name string
	Subscriptions []Subscription
}

func main() {
	// config := ReadYamlFile("sample.yaml")
	// fmt.Printf("%#v",config)
	log.Println("Starting up")
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222")
	defer nc.Drain()
	log.Println("Connected to NATS")
	
	nc.Subscribe("ping", func(m *nats.Msg) {
		logMessageReceived(m)
		actions.Pong(nc)
	})
	nc.Subscribe("ups.office", func(m *nats.Msg) {
		logMessageReceived(m)
		actions.Shutdown(nc)
	})

	log.Println("Subscription setup complete")
	for {
		time.Sleep(1 * time.Second)
	}
}

func logMessageReceived(m *nats.Msg) {
	log.Printf("Received a message from subject %s: %s \n", m.Subject, string(m.Data))
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