package project

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
