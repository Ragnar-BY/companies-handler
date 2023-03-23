package service

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
)

// UserRepository describes user repository
type UserRepository interface {
	CreateUser(ctx context.Context, u domain.User) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates new user
func (s *UserService) CreateUser(ctx context.Context, u domain.User) (*domain.User, error) {
	return s.repo.CreateUser(ctx, u)
}
