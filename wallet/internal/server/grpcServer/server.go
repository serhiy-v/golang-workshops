package grpcserver

import (
	"context"
	"database/sql"
	"log"

	"github.com/workshops/wallet/internal/middleware/auth"
	pb "github.com/workshops/wallet/internal/proto"
	"github.com/workshops/wallet/internal/repository/postgre"
	"github.com/workshops/wallet/internal/services/wallet"
)

// type Validator interface {
//	Validate(interface{}) error
//}

type Server struct {
	// valid Validator
	service    *wallet.Service
	jwtWrapper *auth.JwtWrapper
	pb.UserServiceServer
	pb.WalletServiceServer
	pb.TransactionServiceServer
}

func NewGrpcServer(service *wallet.Service, jwtWrapper *auth.JwtWrapper) *Server {
	return &Server{
		service:    service,
		jwtWrapper: jwtWrapper,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	name := req.GetName()

	token, err := s.jwtWrapper.GenerateToken(name)
	if err != nil {
		log.Printf("Unable to generate token: %v\n", err)
		return nil, err
	}

	err = s.service.CreateUser(token)
	if err != nil {
		log.Printf("Unable to create: %v\n", err)
		return nil, err
	}

	res := &pb.CreateUserResponse{
		Name: name,
	}
	return res, nil
}

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.service.GetUsers()
	if err != nil {
		log.Printf("Unable to get users : %v\n", err)
		return nil, err
	}

	var protoUsers []*pb.User

	for _, user := range users {
		protoUsers = append(protoUsers, convertUser(user))
	}

	res := &pb.GetUsersResponse{
		User: protoUsers,
	}

	return res, nil
}

func (s *Server) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	wallet := &postgre.Wallet{
		Balance: int(req.GetBalance()),
		UserID:  req.GetUserId(),
	}

	err := s.service.CreateWallet(wallet)
	if err != nil {
		log.Printf("Unable to create: %v\n", err)
		return nil, err
	}

	res := &pb.CreateWalletResponse{
		Balance: req.GetBalance(),
		UserId:  req.GetUserId(),
	}
	return res, nil
}

func (s *Server) GetWalletByID(ctx context.Context, req *pb.GetWalledByIdRequest) (*pb.GetWalletByIdResponse, error) {
	id := req.GetId()

	wallet, err := s.service.GetWalletByID(id)
	if err != nil {
		log.Printf("Unable to get wallet: %v\n", err)
		return nil, err
	}

	pbWallet := &pb.Wallet{
		Id:      wallet.ID,
		Balance: int32(wallet.Balance),
		UserId:  wallet.UserID,
	}

	res := &pb.GetWalletByIdResponse{
		Wallet: pbWallet,
	}

	return res, nil
}

func (s *Server) GetTransactions(ctx context.Context,
	req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	transactions, err := s.service.GetTransactions()
	if err != nil {
		log.Printf("Unable to get transactions : %v\n", err)
		return nil, err
	}

	var protoTransactions []*pb.Transaction

	for _, transaction := range transactions {
		protoTransactions = append(protoTransactions, convertTransaction(transaction))
	}

	res := &pb.GetTransactionResponse{
		Transaction: protoTransactions,
	}

	return res, nil
}

func (s *Server) CreateTransaction(ctx context.Context,
	req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	transaction := &postgre.Transaction{
		CreditWalletID: req.GetCreditWalletId(),
		DebitWalletID:  req.GetDebitWalletId(),
		Amount:         int(req.GetAmount()),
	}

	err := s.service.CreateTransaction(transaction)
	if err != nil {
		log.Printf("Transaction Failled: %v\n", err)
		return nil, err
	}

	res := &pb.CreateTransactionResponse{
		CreditWalletId: req.GetCreditWalletId(),
		DebitWalletId:  req.GetDebitWalletId(),
		Amount:         req.GetAmount(),
	}

	return res, nil
}

func (s *Server) GetWalletTransactionsByID(ctx context.Context,
	req *pb.GetWalletTransactionsByIdRequest) (*pb.GetTransactionResponse, error) {
	id := req.GetId()

	transactions, err := s.service.GetWalletTransactionsByID(id)
	if err != nil {
		log.Printf("Unable to get transactions : %v\n", err)
		return nil, err
	}
	var protoTransactions []*pb.Transaction

	for _, transaction := range transactions {
		protoTransactions = append(protoTransactions, convertTransaction(transaction))
	}

	res := &pb.GetTransactionResponse{
		Transaction: protoTransactions,
	}

	return res, nil
}

func convertTransaction(transaction *postgre.Transaction) *pb.Transaction {
	return &pb.Transaction{
		Id:             transaction.ID,
		CreditWalletId: transaction.CreditWalletID,
		DebitWalletId:  transaction.DebitWalletID,
		Amount:         int32(transaction.Amount),
		Type:           int32(transaction.Type),
		FeeAmount:      int32(transaction.FeeAmount),
		FeeWalletId:    transaction.FeeWalletID,
		CreditUserId:   transaction.CreditUserID,
		DebitUserId:    transaction.DebitUserID,
	}
}

func convertUser(user *postgre.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Token: convertNullStr(user.Token),
	}
}

func convertNullStr(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
