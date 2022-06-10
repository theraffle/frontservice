package main

import (
	"context"
	"fmt"
	"github.com/theraffle/frontservice/src/internal/logrotate"
	"github.com/theraffle/frontservice/src/pkg/server"
	"io"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

const (
	port = "8080"
)

func main() {
	ctx := context.Background()
	// Set log rotation
	logFile, err := logrotate.LogFile()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = logFile.Close()
	}()
	logWriter := io.MultiWriter(logFile, os.Stdout)
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(logWriter)))
	if err := logrotate.StartRotate("0 0 1 * * ?"); err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}

	// set port
	srvPort := port
	if os.Getenv("PORT") != "" {
		srvPort = os.Getenv("PORT")
	}
	srv, err := server.New(ctx)
	if err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}
	srv.Start(srvPort)
}
