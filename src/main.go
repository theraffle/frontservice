package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/theraffle/frontservice/src/pkg/server"
	"os"
	"time"
)

const (
	port = "8080"
)

func main() {
	ctx := context.Background()
	// set log
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	// set port
	srvPort := port
	if os.Getenv("PORT") != "" {
		srvPort = os.Getenv("PORT")
	}
	srv, err := server.New(ctx)
	if err != nil {
		log.Errorf(err.Error(), "")
		os.Exit(1)
	}
	srv.Start(srvPort)
}
