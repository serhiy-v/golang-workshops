package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/workshops/wallet/internal/middleware/auth"
	"github.com/workshops/wallet/internal/repository/models"
	"github.com/workshops/wallet/internal/services/wallet"
)

type Validator interface {
	Validate(interface{}) error
}

type Server struct {
	valid          Validator
	jwtWrapper     *auth.JwtWrapper
	servicePostgre *wallet.Service
	serviceMongo   *wallet.Service
}

func NewServer(servicePostgre *wallet.Service, serviceMongo *wallet.Service, jwtWrapper *auth.JwtWrapper, validator Validator) *Server {
	return &Server{
		servicePostgre: servicePostgre,
		serviceMongo:   serviceMongo,
		jwtWrapper:     jwtWrapper,
		valid:          validator,
	}
}

func (s *Server) RunServer() {
	router := NewRouter(s)
	addr := "localhost:8090"

	log.Printf("start HTTP server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	var user models.User
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

	switch db {
	case "mongo":
		err = s.serviceMongo.CreateUser(token)
		if err != nil {
			http.Error(w, "Unable to create user", http.StatusForbidden)
			log.Printf("Unable to create: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created user: " + user.Name))
	case "postgre":
		err = s.servicePostgre.CreateUser(token)
		if err != nil {
			http.Error(w, "Unable to create user", http.StatusForbidden)
			log.Printf("Unable to create: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created user: " + user.Name))
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	switch db {
	case "mongo":
		users, err := s.serviceMongo.GetUsers()
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
	case "postgre":
		users, err := s.servicePostgre.GetUsers()
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
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) CreateWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	var wallet models.Wallet
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		log.Printf("Unable to get wallet from request: %v\n", err)
	}

	err = s.valid.Validate(wallet)
	if err != nil {
		log.Printf("Bad input: %v\n", err)
		return
	}

	switch db {
	case "mongo":
		err = s.serviceMongo.CreateWallet(&wallet)
		if err != nil {
			log.Printf("Unable to create: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wallet)
	case "postgre":
		err = s.servicePostgre.CreateWallet(&wallet)
		if err != nil {
			log.Printf("Unable to create: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(wallet)
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	params = mux.Vars(r)
	id := params["id"]

	switch db {
	case "mongo":
		wallet, err := s.serviceMongo.GetWalletByID(id)
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
	case "postgre":
		wallet, err := s.servicePostgre.GetWalletByID(id)
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
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetWalletTransactionsByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	params = mux.Vars(r)
	id := params["id"]

	switch db {
	case "mongo":
		transactions, err := s.serviceMongo.GetWalletTransactionsByID(id)
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
	case "postgre":
		transactions, err := s.servicePostgre.GetWalletTransactionsByID(id)
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
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	switch db {
	case "mongo":
		transactions, err := s.serviceMongo.GetTransactions()
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
	case "postgre":
		transactions, err := s.servicePostgre.GetTransactions()
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
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) CreateTransactions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		log.Printf("Unable to get transaction from request: %v\n", err)
	}

	err = s.valid.Validate(transaction)
	if err != nil {
		log.Printf("Bad input: %v\n", err)
		return
	}

	switch db {
	case "mongo":
		err = s.serviceMongo.CreateTransaction(&transaction)
		if err != nil {
			http.Error(w, "Unable to create transaction", http.StatusForbidden)
			log.Printf("Transaction Failled: %v\n", err)
			return
		}

		json.NewEncoder(w).Encode(transaction)
	case "postgre":
		err = s.servicePostgre.CreateTransaction(&transaction)
		if err != nil {
			http.Error(w, "Unable to create transaction", http.StatusForbidden)
			log.Printf("Transaction Failled: %v\n", err)
			return
		}

		json.NewEncoder(w).Encode(transaction)
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetWalletAmountDayByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	params = mux.Vars(r)
	id := params["id"]

	var day models.Week
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		log.Printf("Unable to get day from request: %v\n", err)
	}

	switch db {
	case "mongo":
		days, err := s.servicePostgre.GetWalletAmountDayByID(id, day)
		if err != nil {
			http.Error(w, "Unable to get amount per day:", http.StatusForbidden)
			log.Printf("Unable to get amount per day:: %v\n", err)
			return
		}

		w.Write([]byte(fmt.Sprintf("Amount:%v", days)))
	case "postgre":
		days, err := s.servicePostgre.GetWalletAmountDayByID(id, day)
		if err != nil {
			http.Error(w, "Unable to get amount per day", http.StatusForbidden)
			log.Printf("Unable to get amount per day: %v\n", err)
			return
		}

		for _, day := range days {
			err = json.NewEncoder(w).Encode(day)
			if err != nil {
				log.Printf("Unable to encode transaction: %v\n", err)
				return
			}
		}
	default:
		w.Write([]byte("invalid db"))
	}
}

func (s *Server) GetWalletAmountWeekByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := params["db"]

	params = mux.Vars(r)
	id := params["id"]

	var day models.Week
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		log.Printf("Unable to get day from request: %v\n", err)
	}

	switch db {
	case "mongo":
		incomeAmount, outcomeAmount, err := s.serviceMongo.GetWalletAmountWeekByID(id, day)
		if err != nil {
			http.Error(w, "Unable to get amount per day:", http.StatusForbidden)
			log.Printf("Unable to get amount per day:: %v\n", err)
			return
		}

		w.Write([]byte(fmt.Sprintf("Amount:%v", incomeAmount, outcomeAmount)))
	case "postgre":
		incomeAmount, outcomeAmount, err := s.servicePostgre.GetWalletAmountWeekByID(id, day)
		if err != nil {
			http.Error(w, "Unable to get amount per day", http.StatusForbidden)
			log.Printf("Unable to get amount per day: %v\n", err)
			return
		}
		//value1 := incomeAmount
		//value2 := outcomeAmount
		w.Write([]byte(fmt.Sprintf("day:%v, incomeAmount:%v, outcomeAmount:%v", day, incomeAmount, outcomeAmount)))
	default:
		w.Write([]byte("invalid db"))
	}
}
