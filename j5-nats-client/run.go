package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/jacksondr5/go-monorepo/j5-nats-client/actions"
	natscommon "github.com/jacksondr5/go-monorepo/nats-common"
	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

type Subscription struct {
	Action string
	Name string
	Trigger string
	Subject string
}

type Config struct {
	Name string
	Subscriptions []Subscription
}

var hostname string

func (p *program) run() {
	dir, err := os.Getwd()
	if err != nil {
		logger.Info("Error getting current working directory")
		logFatal(logger, err.Error())
	}
	logger.Infof("Current working directory: %s", dir)
	if strings.Contains(dir, "system32") {
		logger.Info("Running on a Windows machine, changing directory to C:\\j5")
		err := os.Chdir("C:\\j5");
		if err != nil {
			logger.Info("Error changing directory to C:\\j5")
			logFatal(logger, err.Error())
		}
	}
	config := ReadYamlFile("j5-nats-client-config.yaml")
	logger.Infof("%#v",config)

	hostname = getHostname()
	logger.Info("Starting up")
	nc, _ := nats.Connect("nats://nats.k8s.j5:4222", nats.Name(hostname))
	defer nc.Drain()
	logger.Info("Connected to NATS")

	logger.Info("Setting up subscriptions")
	for i := range config.Subscriptions {
		subscription := config.Subscriptions[i]
		logger.Infof("Setting up subscription for \"%s\" on subject \"%s\" with action \"%s\" and trigger \"%s\"", subscription.Name, subscription.Subject, subscription.Action, subscription.Trigger)
		nc.Subscribe(subscription.Subject, func(m *nats.Msg) {
			natscommon.LogMessageReceived(m)
			if subscription.Trigger != "" {
				if string(m.Data) != subscription.Trigger {
					logger.Infof("Message \"%s\" on subject \"%s\" does not match trigger \"%s\".  Ignoring.", string(m.Data), m.Subject,  subscription.Trigger)
					return
				} else {
					logger.Infof("Message \"%s\" on subject \"%s\" matches trigger \"%s\"", string(m.Data), m.Subject,  subscription.Trigger)
				}
			} else {
				logger.Infof("No trigger message specified for subscription \"%s\".", subscription.Name)
			}
			logger.Infof("Executing action \"%s\"", subscription.Action)
			switch subscription.Action {
			case "pong":
				actions.Pong(nc, logger, hostname)
			case "shutdown":
				actions.ShutdownUbuntu(nc, logger, hostname, subscription.Subject)
			case "shutdown-windows":
				actions.ShutdownWindows(nc, logger, hostname, subscription.Subject)
			case "shutdown-gitlab":
				actions.ShutdownGitLabProcess(nc, logger, hostname, subscription.Subject)
			default:
				logger.Infof("Unknown action %s", subscription.Action)
			}
		})
	}


	logger.Info("Subscription setup complete.  Waiting for events.")

	http.ListenAndServe(":12345", nil)
}


func getHostname() string {
	hostname, hostnameErr := os.Hostname()
	if hostnameErr != nil {
		logFatal(logger, hostnameErr.Error())
	}
	return hostname
}

func ReadYamlFile(path string) *Config {
	buf, err := os.ReadFile(path)
	if err != nil {
		logFatal(logger, err.Error())
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		logFatal(logger, err.Error())
	}

	return &config
}