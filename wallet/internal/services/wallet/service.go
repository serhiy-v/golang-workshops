package wallet

import "github.com/workshops/wallet/internal/repository/postgre"

type Repository interface {
	CreateUser(token string) error
	CreateWallet(wallet *postgre.Wallet) error
	GetUsers() ([]*postgre.User, error)
	GetWalletByID(id string) (*postgre.Wallet, error)
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

func (s *Service) CreateTransaction(transaction *postgre.Transaction) error {
	return s.repo.CreateTransaction(transaction)
}
