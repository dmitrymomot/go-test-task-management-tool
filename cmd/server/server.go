package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

func newServer(r routerInterface, l logger) *server {
	return &server{r, l}
}

type server struct {
	r routerInterface
	l logger
}

type logger interface {
	Error(err error)
	Info(msg string)
}

type routerInterface interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (s *server) run(addr string) {
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Accept",
			"Content-Type",
			"Origin",
		}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		}),
	)
	srv := &http.Server{
		Handler:      cors(s.r),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.l.Info(fmt.Sprintf("Server started %s", addr))
	s.l.Error(srv.ListenAndServe())
}
