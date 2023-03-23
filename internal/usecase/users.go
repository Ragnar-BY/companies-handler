package usecase

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
)

// UserService describes user service
type UserService interface {
	CreateUser(ctx context.Context, u domain.User) (*domain.User, error)
}

// UserUsecase is user usecase
type UserUsecase struct {
	srv UserService
}

// NewUserUsecase creates new user usecase
func NewUserUsecase(srv UserService) *UserUsecase {
	return &UserUsecase{srv: srv}
}

// CreateUser creates new user
func (s *UserUsecase) CreateUser(ctx context.Context, u domain.User) (*domain.User, error) {
	return s.srv.CreateUser(ctx, u)
}
