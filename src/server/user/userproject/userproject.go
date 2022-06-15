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

package userproject

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
	ctx         context.Context
	log         logr.Logger
	userSvcConn *grpc.ClientConn
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, log logr.Logger, userSvcConn *grpc.ClientConn) (apihandler.APIHandler, error) {
	handler := &handler{ctx: ctx, log: log, userSvcConn: userSvcConn}

	// Create User Project
	createUserProject := wrapper.New("/project", []string{http.MethodPost}, handler.createUserProjectHandler)
	if err := parent.Add(createUserProject); err != nil {
		return nil, err
	}

	// Get User Projects
	getUserProjects := wrapper.New("/projects", []string{http.MethodGet}, handler.getUserProjectsHandler)
	if err := parent.Add(getUserProjects); err != nil {
		return nil, err
	}

	return handler, nil
}

type createUserProjectReqBody struct {
	ProjectID int64  `json:"project_id,omitempty"`
	ChainID   int64  `json:"chain_id,omitempty"`
	Address   string `json:"address,omitempty"`
}

func (h handler) createUserProjectHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("create_user_project_request", reqID)

	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}

	log.Info("create user project", "id", id)

	intID, _ := strconv.Atoi(id)
	// Decode request body
	createUserProjectReq := &createUserProjectReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(createUserProjectReq); err != nil {
		h.log.Error(err, "create user project error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}
	// TODO request validity check

	resp, err := pb.NewUserServiceClient(h.userSvcConn).CreateUserProject(h.ctx, &pb.CreateUserProjectRequest{
		UserID:    int64(intID),
		ProjectID: createUserProjectReq.ProjectID,
		ChainID:   createUserProjectReq.ChainID,
		Address:   createUserProjectReq.Address,
	})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h handler) getUserProjectsHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("get_user_project_request", reqID)

	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}

	log.Info("getting user projects info", "id", id)

	intID, _ := strconv.Atoi(id)
	resp, err := pb.NewUserServiceClient(h.userSvcConn).GetUserProject(h.ctx, &pb.GetUserProjectRequest{UserID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}
