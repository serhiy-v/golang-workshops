package http

import (
	"encoding/json"
	"fmt"
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
	http.ListenAndServe("localhost:8090", router)
}

func (s *Server) HandlerA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
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
	err = s.service.CreateUser(token)
	if err != nil {
		log.Printf("Unable to create: %v\n", err)
		return
	}
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.service.GetUsers()
	if err != nil {
		log.Printf("Unable to get users : %v\n", err)
	}
	for _, user := range users {
		json.NewEncoder(w).Encode(user)
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
}

func (s *Server) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	wallet, err := s.service.GetWalletByID(id)
	if err != nil {
		log.Printf("Unable to get wallet: %v\n", err)
		return
	}
	json.NewEncoder(w).Encode(wallet)
}

func (s *Server) GetWalletTransactionsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	transactions, err := s.service.GetWalletTransactionsById(id)
	if err != nil {
		log.Printf("Unable to get transactions : %v\n", err)
	}
	for _, transaction := range transactions {
		json.NewEncoder(w).Encode(transaction)
	}
}

func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := s.service.GetTransactions()
	if err != nil {
		log.Printf("Unable to get transactions : %v\n", err)
	}
	for _, transaction := range transactions {
		json.NewEncoder(w).Encode(transaction)
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
		log.Printf("Transaction Failled: %v\n", err)
		return
	}

}
