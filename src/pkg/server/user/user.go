package user

import (
	"context"
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	pb "github.com/theraffle/frontservice/src/genproto/pb"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/utils"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"github.com/theraffle/frontservice/src/pkg/server/user/userproject"
	"github.com/theraffle/frontservice/src/pkg/server/user/wallet"
	"google.golang.org/grpc"
	"net/http"
)

type handler struct {
	ctx context.Context
	log logr.Logger

	userSvcAddr    string
	userSvcConn    *grpc.ClientConn
	projectHandler apiserver.APIHandler
	walletHandler  apiserver.APIHandler
}

type createUserReqBody struct {
	UserID    string       `json:"user_id"`
	LoginType pb.LoginType `json:"login_type"`
}

// NewHandler instantiates a new apis handler
func NewHandler(ctx context.Context, parent wrapper.RouterWrapper, logger logr.Logger) (apiserver.APIHandler, error) {
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
	walletHandler, err := wallet.NewHandler(userWrapper, logger)
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
	log := h.log.WithValues("request", reqID)
	id := mux.Vars(req)["id"]
	if id == "" {
		_ = utils.RespondError(w, http.StatusBadRequest, "user id not specified")
		return
	}
	log.Info("getting info of user", "id", id)
	resp, err := pb.NewUserServiceClient(h.userSvcConn).GetUser(h.ctx, &pb.GetUserRequest{UserID: int64(64)})
	if err != nil {
		h.log.Error(err, "")
		_ = utils.RespondError(w, http.StatusBadRequest, "response error")
	}
	_ = utils.RespondJSON(w, resp)
}

func (h *handler) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) getUserProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}
