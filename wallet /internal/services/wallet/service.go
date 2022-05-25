package wallet

type Repository interface {
	CreateUser(token string) error
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
