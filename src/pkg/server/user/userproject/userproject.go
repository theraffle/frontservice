package userproject

import (
	"github.com/go-logr/logr"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"net/http"
)

type handler struct {
	log logr.Logger
}

// NewHandler instantiates a new apis handler
func NewHandler(parent wrapper.RouterWrapper, _ logr.Logger) (apiserver.APIHandler, error) {
	handler := &handler{}

	// Create User Project
	createUserProject := wrapper.New("/userproject", []string{http.MethodPost}, handler.createUserProjectHandler)
	if err := parent.Add(createUserProject); err != nil {
		return nil, err
	}

	// Get User Projects
	getUserProject := wrapper.New("/userproject", []string{http.MethodGet}, handler.getUserProjectHandler)
	if err := parent.Add(getUserProject); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h handler) getUserProjectHandler(writer http.ResponseWriter, request *http.Request) {

}

func (h handler) createUserProjectHandler(writer http.ResponseWriter, request *http.Request) {

}
