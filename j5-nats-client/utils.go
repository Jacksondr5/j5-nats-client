package main

import (
	"os"

	"github.com/kardianos/service"
)

func logFatal(l service.Logger, msg string) {
	l.Error(msg)
	os.Exit(1)
}