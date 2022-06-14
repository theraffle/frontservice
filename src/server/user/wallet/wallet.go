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

package wallet

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

type createUserWalletReqBody struct {
	ChainID int64  `json:"chain_id,omitempty"`
	Address string `json:"address,omitempty"`
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, log logr.Logger, userSvcConn *grpc.ClientConn) (apihandler.APIHandler, error) {
	handler := &handler{ctx: ctx, log: log, userSvcConn: userSvcConn}

	// Create User Wallet
	createUserWallet := wrapper.New("/wallet", []string{http.MethodPost}, handler.createUserWalletHandler)
	if err := parent.Add(createUserWallet); err != nil {
		return nil, err
	}

	// Get User Wallets
	getUserWallet := wrapper.New("/wallets", []string{http.MethodGet}, handler.getUserWalletHandler)
	if err := parent.Add(getUserWallet); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h handler) createUserWalletHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("request", reqID)

	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}

	log.Info("create user wallet", "id", id)

	intID, _ := strconv.Atoi(id)
	// Decode request body
	createUserWalletReq := &createUserWalletReqBody{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(createUserWalletReq); err != nil {
		h.log.Error(err, "create user wallet error")
		_ = utils.RespondError(w, http.StatusBadRequest, "request body is not in json form or is malformed")
		return
	}
	// TODO request validity check

	resp, err := pb.NewUserServiceClient(h.userSvcConn).CreateUserWallet(h.ctx, &pb.CreateUserWalletRequest{
		Wallet: &pb.UserWallet{
			UserID:  int64(intID),
			ChainID: createUserWalletReq.ChainID,
			Address: createUserWalletReq.Address,
		},
	})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h handler) getUserWalletHandler(w http.ResponseWriter, req *http.Request) {
	reqID := utils.RandomString(10)
	log := h.log.WithValues("request", reqID)

	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}

	log.Info("getting user wallet info", "id", id)

	intID, _ := strconv.Atoi(id)
	resp, err := pb.NewUserServiceClient(h.userSvcConn).GetUserWallet(h.ctx, &pb.GetUserWalletRequest{UserID: int64(intID)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}
