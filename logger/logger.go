package logger

import (
	"fmt"
	"maps"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

type BasicLogInfo struct {
    Id *int64 `json:"id,omitempty"`
    // Filename *string `json:"filename,omitempty"`
    Level *string `json:"level,omitempty"`
    Message *string `json:"message,omitempty"`
    Source *string `json:"source,omitempty"`
    Timestamp *int64 `json:"timestamp,omitempty"`
}

type ErrorLog struct {
    // Error log.
    Log *BasicLogInfo `json:"log,omitempty"`
    // Error log.
    Error *error `json:"error,omitempty"`
}

func Init() {
	level := os.Getenv("LOG_LEVEL")
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		logLevel = log.InfoLevel
		ErrorWithFields(
			"Error parsing log level", 
			err, 
			log.Fields{
				"level": level,
			},
		)
	}
	log.SetLevel(logLevel)
	DebugWithFields("Setting log level", log.Fields{
		"envVar": level,
		"logLevel": logLevel,
	})
}

func makeEntry(
	level string, 
	msg string, 
	file string, 
	line int, 
	additionalFields log.Fields,
) *log.Entry {
	fields := log.Fields{
		"level": level,
		"msg": msg,
		"source": fmt.Sprintf("%s:%d", file, line),
		"timestamp": time.Now(),
	}
	maps.Copy(fields, additionalFields)
	return log.WithFields(fields)
}

func Debug(msg string) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("debug", msg, file, line, nil).Debug(msg)
}

func DebugWithFields(msg string, fields log.Fields) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("debug", msg, file, line, fields).Debug(msg)
}

func Info(msg string) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("info", msg, file, line, nil).Info(msg)
}

func InfoWithFields(msg string, fields log.Fields) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("info", msg, file, line, fields).Info(msg)
}

func Warning(msg string) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("warning", msg, file, line, nil).Warning(msg)
}

func WarningWithFields(msg string, fields log.Fields) {
	_, file, line, _ := runtime.Caller(0)
	makeEntry("warning", msg, file, line, fields).Warning(msg)
}

func Error(msg string, err error) {
	_, file, line, _ := runtime.Caller(0)
	errorFields := log.Fields{
		"error": err,
	}
	makeEntry("error", msg, file, line, errorFields).Error(msg)
}

func ErrorWithFields(msg string, err error, fields log.Fields) {
	_, file, line, _ := runtime.Caller(0)
	errorFields := log.Fields{
		"error": err,
	}
	maps.Copy(errorFields, fields)
	makeEntry("error", msg, file, line, errorFields).Error(msg)
}

func Fatal(msg string, err error) {
	_, file, line, _ := runtime.Caller(0)
	errorFields := log.Fields{
		"error": err,
	}
	makeEntry("fatal", msg, file, line, errorFields).Fatal(msg)
}
