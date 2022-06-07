package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/workshops/wallet/internal/middleware/auth"
	"github.com/workshops/wallet/internal/repository/postgre"
	"github.com/workshops/wallet/internal/services/wallet"
)

type Validator interface {
	Validate(interface{}) error
}

type Server struct {
	valid      Validator
	jwtWrapper *auth.JwtWrapper
	service    *wallet.Service
}

func NewServer(service *wallet.Service, jwtWrapper *auth.JwtWrapper, validator Validator) *Server {
	return &Server{
		service:    service,
		jwtWrapper: jwtWrapper,
		valid:      validator,
	}
}

func (s *Server) RunServer() {
	router := NewRouter(s)
	addr := "localhost:8090"

	log.Printf("start HTTP server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user postgre.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Unable to get user from request: %v\n", err)
	}

	err = s.valid.Validate(user)
	if err != nil {
		log.Printf("Bad input: %v\n", err)
		return
	}

	token, err := s.jwtWrapper.GenerateToken(user.Name)
	if err != nil {
		log.Printf("Unable to generate token: %v\n", err)
		return
	}

	err = s.service.CreateUser(token)
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusForbidden)
		log.Printf("Unable to create: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created user: " + user.Name))
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.service.GetUsers()
	if err != nil {
		http.Error(w, "Unable to get users", http.StatusForbidden)
		log.Printf("Unable to get users : %v\n", err)
	}

	for _, user := range users {
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Printf("Unable to encode user: %v\n", err)
			return
		}
	}
}

func (s *Server) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var wallet postgre.Wallet
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		log.Printf("Unable to get wallet from request: %v\n", err)
	}

	err = s.valid.Validate(wallet)
	if err != nil {
		log.Printf("Bad input: %v\n", err)
		return
	}

	err = s.service.CreateWallet(&wallet)
	if err != nil {
		log.Printf("Unable to create: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wallet)
}

func (s *Server) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	wallet, err := s.service.GetWalletByID(id)
	if err != nil {
		http.Error(w, "Unable to get wallet", http.StatusForbidden)
		log.Printf("Unable to get wallet: %v\n", err)
		return
	}

	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		log.Printf("Unable to encode wallet: %v\n", err)
		return
	}
}

func (s *Server) GetWalletTransactionsByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	transactions, err := s.service.GetWalletTransactionsByID(id)
	if err != nil {
		http.Error(w, "Unable to get wallet transactions", http.StatusForbidden)
		log.Printf("Unable to get transactions : %v\n", err)
	}

	for _, transaction := range transactions {
		err = json.NewEncoder(w).Encode(transaction)
		if err != nil {
			log.Printf("Unable to encode transaction: %v\n", err)
			return
		}
	}
}

func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := s.service.GetTransactions()
	if err != nil {
		http.Error(w, "Unable to get transactions", http.StatusForbidden)
		log.Printf("Unable to get transactions : %v\n", err)
	}

	for _, transaction := range transactions {
		err = json.NewEncoder(w).Encode(transaction)
		if err != nil {
			log.Printf("Unable to encode transaction: %v\n", err)
			return
		}
	}
}

func (s *Server) CreateTransactions(w http.ResponseWriter, r *http.Request) {
	var transaction postgre.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Printf("Unable to get transaction from request: %v\n", err)
	}

	err = s.valid.Validate(transaction)
	if err != nil {
		log.Printf("Bad input: %v\n", err)
		return
	}

	err = s.service.CreateTransaction(&transaction)
	if err != nil {
		http.Error(w, "Unable to create transaction", http.StatusForbidden)
		log.Printf("Transaction Failled: %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}
