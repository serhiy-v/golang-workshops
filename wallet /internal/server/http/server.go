package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/workshops/wallet/internal/middleware/auth"
	"github.com/workshops/wallet/internal/repository/postgre"
	"github.com/workshops/wallet/internal/services/wallet"
)

//type Validator interface {
//	Validate(interface{}) error
//}

type Server struct {
	//valid Validator
	jwtWrapper *auth.JwtWrapper
	service    *wallet.Service
}

func NewServer(service *wallet.Service, jwtWrapper *auth.JwtWrapper) *Server {
	return &Server{
		service:    service,
		jwtWrapper: jwtWrapper,
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
		log.Println("Unable to get user from request")
	}
	token, err := s.jwtWrapper.GenerateToken(user.Name)
	err = s.service.CreateUser(token)
	if err != nil {
		log.Println("Unable to create")
		return
	}
}
