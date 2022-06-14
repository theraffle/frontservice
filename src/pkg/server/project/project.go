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

package project

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/utils"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"google.golang.org/grpc"
	"net/http"
)

type handler struct {
	ctx context.Context
	_   logr.Logger

	projectSvcAddr string
	projectSvcConn *grpc.ClientConn
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, _ logr.Logger) (apiserver.APIHandler, error) {
	handler := &handler{ctx: ctx}
	utils.MustMapEnv(&handler.projectSvcAddr, "PROJECT_SERVICE_ADDR")
	utils.MustConnGRPC(ctx, &handler.projectSvcConn, handler.projectSvcAddr)

	// Create Project
	createProject := wrapper.New("/project", []string{http.MethodPost}, handler.createProjectHandler)
	if err := parent.Add(createProject); err != nil {
		return nil, err
	}

	// Get All Projects
	getAllProject := wrapper.New("/project", []string{http.MethodGet}, handler.getAllProjectHandler)
	if err := parent.Add(getAllProject); err != nil {
		return nil, err
	}

	// Get Certain Project
	getProject := wrapper.New("/project/{id}", []string{http.MethodGet}, handler.getProjectHandler)
	if err := parent.Add(getProject); err != nil {
		return nil, err
	}
	// Edit Project
	updateProject := wrapper.New("/project/{id}", []string{http.MethodPut}, handler.updateProjectHandler)
	if err := parent.Add(updateProject); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h *handler) createProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) getAllProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) updateProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) getProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}
