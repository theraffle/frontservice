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
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/theraffle/frontservice/src/apihandler"
	"github.com/theraffle/frontservice/src/genproto/pb"
	"github.com/theraffle/frontservice/src/utils"
	"github.com/theraffle/frontservice/src/wrapper"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type handler struct {
	ctx context.Context
	log logr.Logger

	projectSvcAddr string
	projectSvcConn *grpc.ClientConn
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, logger logr.Logger) (apihandler.APIHandler, error) {
	handler := &handler{ctx: ctx, log: logger}
	utils.MustMapEnv(&handler.projectSvcAddr, "PROJECT_SERVICE_ADDR")
	utils.MustConnGRPC(ctx, &handler.projectSvcConn, handler.projectSvcAddr)

	// Create Project
	createProject := wrapper.New("/project", []string{http.MethodPost}, handler.createProjectHandler)
	if err := parent.Add(createProject); err != nil {
		return nil, err
	}

	// Get All Projects
	getAllProject := wrapper.New("/projects", []string{http.MethodGet}, handler.getAllProjectHandler)
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

type createProjectReqBody struct {
	ProjectName    string `json:"project_name,omitempty"`
	ChainID        int64  `json:"chain_id,omitempty"`
	RaffleContract string `json:"raffle_contract,omitempty"`
}

func (h *handler) createProjectHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("create_project_request", reqID)

	log.Info("create project request")
	// Decode request body
	createProjectReq := &createProjectReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(createProjectReq); err != nil {
		h.log.Error(err, "create project error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}
	// TODO request validity check

	resp, err := pb.NewProjectServiceClient(h.projectSvcConn).CreateProject(h.ctx, &pb.CreateProjectRequest{
		ProjectName:    createProjectReq.ProjectName,
		ChainID:        createProjectReq.ChainID,
		RaffleContract: createProjectReq.RaffleContract,
	})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h *handler) getProjectHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("get_project_request", reqID)
	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "project id not specified")
		return
	}
	log.Info("getting user info", "id", id)
	intID, _ := strconv.Atoi(id)
	resp, err := pb.NewProjectServiceClient(h.projectSvcConn).GetProject(h.ctx, &pb.GetProjectRequest{ProjectID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h *handler) getAllProjectHandler(w http.ResponseWriter, _ *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("get_all_project_request", reqID)

	log.Info("getting all projects list")

	resp, err := pb.NewProjectServiceClient(h.projectSvcConn).GetAllProjects(h.ctx, &pb.Empty{})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

type updateProjectReqBody struct {
}

func (h *handler) updateProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: modify when project components are decided
	reqID := utils.RandomString(10)
	log := h.log.WithValues("update_project_request", reqID)
	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "project id not specified")
		return
	}
	// Decode request body
	updateProjectReq := &updateProjectReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(updateProjectReq); err != nil {
		h.log.Error(err, "update project error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}

	log.Info("updating project info", "id", id)

	intID, _ := strconv.Atoi(id)
	resp, err := pb.NewProjectServiceClient(h.projectSvcConn).UpdateProject(h.ctx, &pb.UpdateProjectRequest{ProjectID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "update user error")
	}
	_ = utils.RespondJSON(w, resp)
}
