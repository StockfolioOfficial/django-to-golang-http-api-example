package core

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

const (
	Debug = true
)

var DefaultValidator = validator.New()

func init() {
	var formatter log.Formatter = &log.JSONFormatter{}
	var level = log.ErrorLevel
	if Debug {
		formatter = &log.TextFormatter{}
		level = log.DebugLevel
	}

	log.SetFormatter(formatter)
	log.SetReportCaller(true)
	log.SetLevel(level)
}