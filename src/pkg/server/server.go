/*
 Copyright 2022 The Raffle Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/utils"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"github.com/theraffle/frontservice/src/pkg/server/user"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Server is an interface of server
type Server interface {
	Start(string)
}

var (
	log = logf.Log.WithName("front-service")
)

type frontendServer struct {
	wrapper     wrapper.RouterWrapper
	userHandler apiserver.APIHandler
	// projectHandler apiserver.APIHandler
}

func New(ctx context.Context) (Server, error) {
	server := new(frontendServer)
	server.wrapper = wrapper.New("/", nil, server.rootHandler)

	server.wrapper.SetRouter(mux.NewRouter())
	server.wrapper.Router().HandleFunc("/", server.rootHandler)

	// Set apisHandler
	userHandler, err := user.NewHandler(ctx, server.wrapper, log)
	if err != nil {
		return nil, err
	}
	server.userHandler = userHandler

	//projectHandler, err := project.NewHandler(ctx, server.wrapper, log)
	//if err != nil {
	//	return nil, err
	//}
	//server.projectHandler = projectHandler

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
