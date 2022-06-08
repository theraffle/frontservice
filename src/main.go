package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"time"
)

const (
	port = "8080"
)

type frontendServer struct {
	userSvcAddr string
	userSvcConn *grpc.ClientConn

	projectSvcAddr string
	projectSvcConn *grpc.ClientConn
}

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
	addr := fmt.Sprintf("0.0.0.0:%s", srvPort)
	svc := new(frontendServer)
	mustMapEnv(&svc.userSvcAddr, "USER_SERVICE_ADDR")
	mustMapEnv(&svc.projectSvcAddr, "PROJECT_SERVICE_ADDR")

	mustConnGRPC(ctx, &svc.userSvcConn, svc.userSvcAddr)
	mustConnGRPC(ctx, &svc.projectSvcConn, svc.projectSvcAddr)

	r := mux.NewRouter()

	var handler http.Handler = r
	//handler = &logHandler{log: log, next: handler} // add logging
	//handler = ensureSessionID(handler)             // add session ID
	handler = &ochttp.Handler{ // add opencensus instrumentation
		Handler:     handler,
		Propagation: &b3.HTTPFormat{}}

	log.Infof("starting server on " + addr + ":" + srvPort)
	log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
	var err error
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	if err != nil {
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}
