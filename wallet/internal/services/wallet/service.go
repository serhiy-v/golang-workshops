package wallet

import (
	"github.com/gammazero/deque"
	"github.com/workshops/wallet/internal/repository/postgre"
)

type Repository interface {
	CreateUser(token string) error
	CreateWallet(wallet *postgre.Wallet) error
	GetUsers() ([]*postgre.User, error)
	GetWalletByID(id string) (*postgre.Wallet, error)
	GetWalletTransactionsById(id string) ([]*postgre.Transaction, error)
	GetTransactions() ([]*postgre.Transaction, error)
	CreateTransaction(transaction *postgre.Transaction) error
}

// Service holds calendar business logic and works with repository
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(token string) error {
	return s.repo.CreateUser(token)
}

func (s *Service) GetUsers() ([]*postgre.User, error) {
	return s.repo.GetUsers()
}

func (s *Service) CreateWallet(wallet *postgre.Wallet) error {
	return s.repo.CreateWallet(wallet)
}

func (s *Service) GetWalletByID(id string) (*postgre.Wallet, error) {
	return s.repo.GetWalletByID(id)
}

func (s *Service) GetWalletTransactionsById(id string) ([]*postgre.Transaction, error) {
	return s.repo.GetWalletTransactionsById(id)
}

func (s *Service) GetTransactions() ([]*postgre.Transaction, error) {
	return s.repo.GetTransactions()
}

func (s *Service) CreateTransaction(transaction *postgre.Transaction) error {
	q := deque.New()
	if transaction.Type == 1 {
		q.PushFront(transaction)
	} else {
		q.PushBack(transaction)
	}
	for q.Len() != 0 {
		//Consume deque and run transactions in another func
		q.PopFront()
	}
	return s.repo.CreateTransaction(transaction)
}
