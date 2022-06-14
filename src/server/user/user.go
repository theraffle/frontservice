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

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/theraffle/frontservice/src/apihandler"
	"github.com/theraffle/frontservice/src/genproto/pb"
	"github.com/theraffle/frontservice/src/server/user/userproject"
	"github.com/theraffle/frontservice/src/server/user/wallet"
	"github.com/theraffle/frontservice/src/utils"
	"github.com/theraffle/frontservice/src/wrapper"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type handler struct {
	ctx context.Context
	log logr.Logger

	userSvcAddr    string
	userSvcConn    *grpc.ClientConn
	projectHandler apihandler.APIHandler
	walletHandler  apihandler.APIHandler
}

type createUserReqBody struct {
	UserID    string       `json:"user_id"`
	LoginType pb.LoginType `json:"login_type"`
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, logger logr.Logger) (apihandler.APIHandler, error) {
	handler := &handler{ctx: ctx, log: logger}
	utils.MustMapEnv(&handler.userSvcAddr, "USER_SERVICE_ADDR")
	utils.MustConnGRPC(ctx, &handler.userSvcConn, handler.userSvcAddr)

	// Create User & Login
	createUser := wrapper.New("/user", []string{http.MethodPost}, handler.createUserHandler)
	if err := parent.Add(createUser); err != nil {
		return nil, err
	}

	// Get User
	getUser := wrapper.New("/user/{id}", []string{http.MethodGet}, handler.getUserHandler)
	if err := parent.Add(getUser); err != nil {
		return nil, err
	}

	// Edit User
	updateUser := wrapper.New("/user/{id}", []string{http.MethodPut}, handler.updateUserHandler)
	if err := parent.Add(updateUser); err != nil {
		return nil, err
	}

	userWrapper := wrapper.New("/user/{id}", nil, nil)
	if err := parent.Add(userWrapper); err != nil {
		return nil, err
	}

	// /user/{id}/project
	projectHandler, err := userproject.NewHandler(userWrapper, logger)
	if err != nil {
		return nil, err
	}
	handler.projectHandler = projectHandler

	// /user/{id}/wallet
	walletHandler, err := wallet.NewHandler(ctx, userWrapper, logger, handler.userSvcConn)
	if err != nil {
		return nil, err
	}
	handler.walletHandler = walletHandler

	return handler, nil
}

func (h *handler) createUserHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("request", reqID)

	log.Info("create user request")
	// Decode request body
	createUserReq := &createUserReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(createUserReq); err != nil {
		h.log.Error(err, "create user error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}
	// TODO request validity check

	resp, err := pb.NewUserServiceClient(h.userSvcConn).LoginUser(h.ctx, &pb.LoginUserRequest{UserID: createUserReq.UserID, LoginType: createUserReq.LoginType})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h *handler) getUserHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("get_user_request", reqID)
	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}
	log.Info("getting user info", "id", id)
	intID, _ := strconv.Atoi(id)
	resp, err := pb.NewUserServiceClient(h.userSvcConn).GetUser(h.ctx, &pb.GetUserRequest{UserID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h *handler) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("get_user_request", reqID)
	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}
	// Decode request body
	updateUserReq := &createUserReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(updateUserReq); err != nil {
		h.log.Error(err, "create user error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}

	log.Info("updating user info", "id", id)

	intID, _ := strconv.Atoi(id)
	userSvcCli := pb.NewUserServiceClient(h.userSvcConn)
	resp, err := userSvcCli.GetUser(h.ctx, &pb.GetUserRequest{UserID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}

	rpcReq := &pb.UpdateUserRequest{
		UserID:     resp.UserID,
		TelegramID: resp.TelegramID,
		DiscordID:  resp.DiscordID,
		TwitterID:  resp.TwitterID,
	}

	if updateUserReq.LoginType == pb.LoginType_DISCORD {
		rpcReq.DiscordID = updateUserReq.UserID
	} else if updateUserReq.LoginType == pb.LoginType_TELEGRAM {
		rpcReq.TelegramID = updateUserReq.UserID
	} else if updateUserReq.LoginType == pb.LoginType_TWITTER {
		rpcReq.TwitterID = updateUserReq.UserID
	} else {
		err = fmt.Errorf("invalid id type")
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "invalid id type")
	}

	resp, err = userSvcCli.UpdateUser(h.ctx, rpcReq)
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "update user error")
	}
	_ = utils.RespondJSON(w, resp)
}
