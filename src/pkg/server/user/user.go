package user

import (
	"github.com/go-logr/logr"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"github.com/theraffle/frontservice/src/pkg/server/user/project"
	"github.com/theraffle/frontservice/src/pkg/server/user/wallet"
	"net/http"
)

type handler struct {
	log            logr.Logger
	projectHandler apiserver.APIHandler
	walletHandler  apiserver.APIHandler
}

// NewHandler instantiates a new apis handler
func NewHandler(parent wrapper.RouterWrapper, logger logr.Logger) (apiserver.APIHandler, error) {
	handler := &handler{}

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
	projectHandler, err := project.NewHandler(userWrapper, logger)
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
	// TODO: implement here
}

func (h *handler) getUserHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}

func (h *handler) getUserProjectHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: implement here
}
