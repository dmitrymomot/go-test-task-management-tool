package webservices

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dmitrymomot/go-test-task-management-tool/app/domain"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/govalidator"
)

// Predefined errors
var (
	ErrMalformedJSON = errors.New("malformed json")
)

// NewTasks is a factory function,
// returns an instance of Tasks structure
func NewTasks(i taskInteractor, l logger) *Tasks {
	return &Tasks{i, l}
}

type logger interface {
	Error(err error)
	Debug(data interface{})
}

type taskInteractor interface {
	GetAll() ([]domain.Task, error)
	GetNew() ([]domain.Task, error)
	GetCompleted() ([]domain.Task, error)
	Store(title, description string) (domain.Task, error)
	Update(id int64, title, description, status string) (domain.Task, error)
	Complete(id int64) (domain.Task, error)
	Delete(id int64) error
}

// Tasks is a task manager web service structure
type Tasks struct {
	i taskInteractor
	l logger
}

// GetAll function is a http handler
func (ws *Tasks) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := ws.i.GetAll()
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
		}
		ws.response(w, http.StatusOK, response{Data: tasks})
	}
}

// GetNew function is a http handler
func (ws *Tasks) GetNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := ws.i.GetNew()
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
		}
		ws.response(w, http.StatusOK, response{Data: tasks})
	}
}

// GetCompleted function is a http handler
func (ws *Tasks) GetCompleted() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := ws.i.GetCompleted()
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
		}
		ws.response(w, http.StatusOK, response{Data: tasks})
	}
}

// Create function is a http handler
func (ws *Tasks) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}{}
		v := govalidator.New(govalidator.Options{
			Request: r,
			Data:    &form,
			Rules: govalidator.MapData{
				"title":       []string{"required", "min:2", "max:250"},
				"description": []string{"required", "min:5", "max:1000"},
			},
		})
		if err := v.ValidateJSON(); len(err) > 0 {
			if _, ok := err["_error"]; ok {
				ws.response(w, http.StatusBadRequest, response{Error: ErrMalformedJSON.Error()})
				return
			}
			ws.response(w, http.StatusBadRequest, response{Error: err})
			return
		}
		task, err := ws.i.Store(form.Title, form.Description)
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
			return
		}
		ws.response(w, http.StatusOK, response{Data: task})
	}
}

// Update function is a http handler
func (ws *Tasks) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := struct {
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status,omitempty"`
		}{}
		v := govalidator.New(govalidator.Options{
			Request: r,
			Data:    &form,
			Rules: govalidator.MapData{
				"title":       []string{"min:2", "max:250"},
				"description": []string{"min:5", "max:1000"},
				"status":      []string{"in:new,done"},
			},
		})
		if err := v.ValidateJSON(); len(err) > 0 {
			if _, ok := err["_error"]; ok {
				ws.response(w, http.StatusBadRequest, response{Error: ErrMalformedJSON.Error()})
				return
			}
			ws.response(w, http.StatusBadRequest, response{Error: err})
			return
		}

		id := ws.getID(r)
		task, err := ws.i.Update(id, form.Title, form.Description, form.Status)
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
			return
		}
		ws.response(w, http.StatusOK, response{Data: task})
	}
}

// Complete function is a http handler
func (ws *Tasks) Complete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ws.getID(r)
		task, err := ws.i.Complete(id)
		if err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
			return
		}
		ws.response(w, http.StatusOK, response{Data: task})
	}
}

// Delete function is a http handler
func (ws *Tasks) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ws.getID(r)
		if err := ws.i.Delete(id); err != nil {
			ws.l.Error(err)
			ws.response(w, http.StatusInternalServerError, response{Error: err.Error()})
			return
		}
		ws.response(w, http.StatusOK, response{Data: true})
	}
}

func (ws *Tasks) response(w http.ResponseWriter, status int, resp response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (ws *Tasks) getID(r *http.Request) int64 {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.ParseInt(idString, 10, 0)
	return id
}

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}
