package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	// ContentTypeHeader is a content-type header key
	ContentTypeHeader = "Content-Type"
	// ContentTypeJSON is json content type
	ContentTypeJSON = "application/json"
)

func newRouter(ws taskWebservice) *router {
	r := mux.NewRouter()
	return &router{r, ws}
}

type taskWebservice interface {
	GetAll() http.HandlerFunc
	GetNew() http.HandlerFunc
	GetCompleted() http.HandlerFunc
	Create() http.HandlerFunc
	Update() http.HandlerFunc
	Complete() http.HandlerFunc
	Delete() http.HandlerFunc
}

type router struct {
	r  *mux.Router
	ws taskWebservice
}

func (r *router) setup() *router {

	r.r.HandleFunc("/", r.mainPageHandler()).Methods(http.MethodGet)
	r.r.HandleFunc("/health", r.healthCheckHandler).Methods(http.MethodGet)

	r.r.HandleFunc("/api/tasks", r.ws.GetAll()).
		Methods(http.MethodGet).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks/new", r.ws.GetNew()).
		Methods(http.MethodGet).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks/completed", r.ws.GetCompleted()).
		Methods(http.MethodGet).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks", r.ws.Create()).
		Methods(http.MethodPost).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks/{id:[0-9]+}", r.ws.Update()).
		Methods(http.MethodPut).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks/{id:[0-9]+}/done", r.ws.Complete()).
		Methods(http.MethodPatch).
		Headers(ContentTypeHeader, ContentTypeJSON)

	r.r.HandleFunc("/api/tasks/{id:[0-9]+}", r.ws.Delete()).
		Methods(http.MethodDelete).
		Headers(ContentTypeHeader, ContentTypeJSON)

	return r
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

func (r *router) mainPageHandler() http.HandlerFunc {
	tpl := template.Must(template.ParseFiles("/tpl/main.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func (r *router) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      http.StatusText(http.StatusOK),
		"server_time": time.Now(),
	})
}

func (r *router) notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": http.StatusText(http.StatusNotFound),
	})
}
