package main

import (
	"log"
	"net"

	"github.com/workshops/wallet/internal/middleware/auth"
	pb "github.com/workshops/wallet/internal/proto"
	"github.com/workshops/wallet/internal/repository/mongo"
	"github.com/workshops/wallet/internal/repository/postgre"
	grpcserver "github.com/workshops/wallet/internal/server/grpcServer"
	"github.com/workshops/wallet/internal/server/http"
	"github.com/workshops/wallet/internal/services/validator"
	"github.com/workshops/wallet/internal/services/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	//"github.com/workshops/wallet/internal/server/http"
)

func main() {
	// main server code
	runHTTP()
	// runGrpc()
}

func runHTTP() {
	str := "postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable"
	mongoStr := "mongodb://gouser:gopassword@localhost:27017/"

	db, err := postgre.NewPostgresDB(str)
	if err != nil {
		log.Fatal(err)
	}

	dbm, err := mongo.NewMongoDB(mongoStr)
	if err != nil {
		log.Fatal(err)
	}

	repoMongo := mongo.NewRepository(dbm)

	repoPostgre := postgre.NewRepository(db)

	validate := validator.NewValidator()
	servicePostgre := wallet.NewService(repoPostgre)
	serviceMongo := wallet.NewService(repoMongo)
	wrapper := auth.NewJwtWrapper("verysecretkey", 999)
	server := http.NewServer(servicePostgre, serviceMongo, wrapper, validate)

	server.RunServer()
}

//nolint
func runGrpc() {
	str := "postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable"

	db, err := postgre.NewPostgresDB(str)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgre.NewRepository(db)
	service := wallet.NewService(repo)
	wrapper := auth.NewJwtWrapper("verysecretkey", 999)
	interceptor := grpcserver.NewAuthInterceptor(wrapper)

	srv := grpcserver.NewGrpcServer(service, wrapper)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
	pb.RegisterUserServiceServer(grpcServer, srv)
	pb.RegisterWalletServiceServer(grpcServer, srv)
	pb.RegisterTransactionServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}

}
