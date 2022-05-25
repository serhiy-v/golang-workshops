package http

import (
	"github.com/gorilla/mux"
)

// will hold http routes and will registrate them
func NewRouter(s *Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", s.CreateUser).Methods("POST")
	sec := r.PathPrefix("/abc").Subrouter()
	sec.Use(s.jwtWrapper.AuthMiddleware)
	sec.HandleFunc("/", s.HandlerA).Methods("GET")

	return r
}
