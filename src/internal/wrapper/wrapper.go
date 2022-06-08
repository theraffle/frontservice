package wrapper

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// RouterWrapper is an interface for wrapper
type RouterWrapper interface {
	Add(child RouterWrapper) error
	FullPath() string

	Router() *mux.Router
	SetRouter(*mux.Router)

	Children() []RouterWrapper

	Parent() RouterWrapper
	SetParent(RouterWrapper)

	Handler() http.HandlerFunc
	SubPath() string
	Methods() []string
}

// Wrapper wraps router with tree structure
type Wrapper struct {
	router *mux.Router

	subPath string
	methods []string
	handler http.HandlerFunc

	children []RouterWrapper
	parent   RouterWrapper
}

// New is a constructor for the wrapper
func New(path string, methods []string, handler http.HandlerFunc) *Wrapper {
	return &Wrapper{
		subPath: path,
		methods: methods,
		handler: handler,
	}
}

// Router returns its router
func (w *Wrapper) Router() *mux.Router {
	return w.router
}

// SetRouter sets its router
func (w *Wrapper) SetRouter(r *mux.Router) {
	w.router = r
}

// Children returns its children
func (w *Wrapper) Children() []RouterWrapper {
	return w.children
}

// Parent returns its parent
func (w *Wrapper) Parent() RouterWrapper {
	return w.parent
}

// SetParent sets parent
func (w *Wrapper) SetParent(parent RouterWrapper) {
	w.parent = parent
}

// Handler returns its handler
func (w *Wrapper) Handler() http.HandlerFunc {
	return w.handler
}

// SubPath returns its subPath
func (w *Wrapper) SubPath() string {
	return w.subPath
}

// Methods returns its methods
func (w *Wrapper) Methods() []string {
	return w.methods
}

// Add adds child as a child (child node of a tree) of w
func (w *Wrapper) Add(child RouterWrapper) error {
	if child == nil || child.(*Wrapper) == nil {
		return fmt.Errorf("child is nil")
	}

	if child.Parent() != nil {
		return fmt.Errorf("child already has parent")
	}

	if child.SubPath() == "" || child.SubPath() == "/" || child.SubPath()[0] != '/' {
		return fmt.Errorf("child subpath is not valid")
	}

	if w.router == nil {
		return fmt.Errorf("parent does not have a router")
	}

	child.SetParent(w)
	w.children = append(w.children, child)

	child.SetRouter(w.router.PathPrefix(child.SubPath()).Subrouter())

	if child.Handler() != nil {
		if len(child.Methods()) > 0 {
			child.Router().Methods(child.Methods()...).Subrouter().HandleFunc("/", child.Handler())
			w.router.Methods(child.Methods()...).Subrouter().HandleFunc(child.SubPath(), child.Handler())
		} else {
			child.Router().HandleFunc("/", child.Handler())
			w.router.HandleFunc(child.SubPath(), child.Handler())
		}
	}

	return nil
}

// FullPath builds full path string of the api
func (w *Wrapper) FullPath() string {
	if w.parent == nil {
		return w.subPath
	}
	re := regexp.MustCompile(`/{2,}`)
	return re.ReplaceAllString(w.parent.FullPath()+w.subPath, "/")
}
