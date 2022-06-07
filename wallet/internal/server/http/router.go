package http

import (
	"github.com/gorilla/mux"
)

// will hold http routes and will registrate them.
func NewRouter(s *Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", s.GetUsers).Methods("GET")
	r.HandleFunc("/users", s.CreateUser).Methods("POST")

	sec := r.PathPrefix("/wallets").Subrouter()
	sec.Use(s.jwtWrapper.AuthMiddleware)
	sec.HandleFunc("", s.CreateWallet).Methods("POST")
	sec.HandleFunc("/{id}", s.GetWalletByID).Methods("GET")
	sec.HandleFunc("/{id}/transactions", s.GetWalletTransactionsByID).Methods("GET")

	trn := r.PathPrefix("/transactions").Subrouter()
	trn.Use(s.jwtWrapper.AuthMiddleware)
	trn.HandleFunc("", s.GetTransactions).Methods("GET")
	trn.HandleFunc("", s.CreateTransactions).Methods("PUT")

	return r
}
