package main

import (
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	logger.Info("Starting up")
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	logger.Info("Shutting down")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "j5-nats-client",
		DisplayName: "J5 NATS Client",
		Description: "This is a service that listens on NATS topics and responds to them.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logFatal(logger, err.Error())
	}
	logger, err = s.Logger(nil)
	if err != nil {
		logFatal(logger, err.Error())
	}
	err = s.Run()
	if err != nil {
		logFatal(logger, err.Error())
	}
}