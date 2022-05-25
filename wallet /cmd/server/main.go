package main

import (
	"fmt"

	"github.com/workshops/wallet/internal/middleware/auth"
	"github.com/workshops/wallet/internal/repository/postgre"
	"github.com/workshops/wallet/internal/server/http"
	"github.com/workshops/wallet/internal/services/wallet"
	//"github.com/workshops/wallet/internal/server/http"
)

func main() {
	// main server code
	str := "postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable"
	repo, err := postgre.NewRepository(str)
	if err != nil {
		fmt.Println(err)
	}
	service := wallet.NewService(repo)
	wrapper := auth.NewJwtWrapper("verysecretkey", 24)

	server := http.NewServer(service, wrapper)

	server.RunServer()

	//ser := http.NewServer()

}
