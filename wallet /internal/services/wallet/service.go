package wallet

import (
	"github.com/gammazero/deque"
	"github.com/workshops/wallet/internal/repository/models"
)

type Repository interface {
	CreateUser(token string) error
	CreateWallet(wallet *models.Wallet) error
	GetUsers() ([]*models.User, error)
	GetWalletByID(id string) (*models.Wallet, error)
	GetWalletTransactionsByID(id string) ([]*models.Transaction, error)
	GetTransactions() ([]*models.Transaction, error)
	CreateTransaction(transaction *models.Transaction) error
	GetWalletAmountDayByID(id string, day models.Day) (int, int, error)
	GetWalletAmountWeekByID(id string, day models.Week) (int, int, error)
}

// Service holds calendar business logic and works with repository.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(token string) error {
	return s.repo.CreateUser(token)
}

func (s *Service) GetUsers() ([]*models.User, error) {
	return s.repo.GetUsers()
}

func (s *Service) CreateWallet(wallet *models.Wallet) error {
	return s.repo.CreateWallet(wallet)
}

func (s *Service) GetWalletByID(id string) (*models.Wallet, error) {
	return s.repo.GetWalletByID(id)
}

func (s *Service) GetWalletTransactionsByID(id string) ([]*models.Transaction, error) {
	return s.repo.GetWalletTransactionsByID(id)
}

func (s *Service) GetTransactions() ([]*models.Transaction, error) {
	return s.repo.GetTransactions()
}

func (s *Service) CreateTransaction(transaction *models.Transaction) error {
	q := deque.New()
	if transaction.Type == 1 {
		q.PushFront(transaction)
	} else {
		q.PushBack(transaction)
	}
	for q.Len() != 0 {
		// Consume deque and run transactions in another func
		q.PopFront()
	}
	return s.repo.CreateTransaction(transaction)
}

func (s *Service) GetWalletAmountDayByID(id string, day models.Day) (int, int, error) {
	return s.repo.GetWalletAmountDayByID(id, day)
}

func (s *Service) GetWalletAmountWeekByID(id string, day models.Week) (int, int, error) {
	return s.repo.GetWalletAmountWeekByID(id, day)
}
