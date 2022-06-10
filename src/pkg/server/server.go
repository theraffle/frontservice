package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/utils"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"github.com/theraffle/frontservice/src/pkg/server/project"
	"github.com/theraffle/frontservice/src/pkg/server/user"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// Server is an interface of server
type Server interface {
	Start(string)
}

var (
	log = logf.Log.WithName("user-service")
)

type frontendServer struct {
	wrapper        wrapper.RouterWrapper
	userHandler    apiserver.APIHandler
	projectHandler apiserver.APIHandler

	userSvcAddr string
	userSvcConn *grpc.ClientConn

	projectSvcAddr string
	projectSvcConn *grpc.ClientConn
}

func New(ctx context.Context) (Server, error) {
	server := new(frontendServer)
	mustMapEnv(&server.userSvcAddr, "USER_SERVICE_ADDR")
	mustMapEnv(&server.projectSvcAddr, "PROJECT_SERVICE_ADDR")

	mustConnGRPC(ctx, &server.userSvcConn, server.userSvcAddr)
	mustConnGRPC(ctx, &server.projectSvcConn, server.projectSvcAddr)

	server.wrapper = wrapper.New("/", nil, server.rootHandler)

	server.wrapper.SetRouter(mux.NewRouter())
	server.wrapper.Router().HandleFunc("/", server.rootHandler)

	// Set apisHandler
	userHandler, err := user.NewHandler(server.wrapper, log)
	if err != nil {
		return nil, err
	}
	server.userHandler = userHandler

	projectHandler, err := project.NewHandler(server.wrapper, log)
	if err != nil {
		return nil, err
	}
	server.projectHandler = projectHandler

	return server, nil
}

func (s *frontendServer) Start(port string) {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Info(fmt.Sprintf("Server is running on %s", addr))
	if err := http.ListenAndServe(addr, s.wrapper.Router()); err != nil {
		log.Error(err, "cannot launch http server")
		os.Exit(1)
	}
}

func (s *frontendServer) rootHandler(w http.ResponseWriter, _ *http.Request) {
	paths := metav1.RootPaths{}
	addPath(&paths.Paths, s.wrapper)

	_ = utils.RespondJSON(w, paths)
}

// addPath adds all the leaf API endpoints
func addPath(paths *[]string, w wrapper.RouterWrapper) {
	if w.Handler() != nil {
		*paths = append(*paths, w.FullPath())
	}

	for _, c := range w.Children() {
		addPath(paths, c)
	}
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
