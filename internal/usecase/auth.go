package usecase

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
)

type AuthService interface {
	GenerateJWT(email string, username string) (string, error)
	ValidateToken(signedToken string) error
}

type AuthUsecase struct {
	auth  AuthService
	users UserService
}

func NewAuthUsecase(auth AuthService, users UserService) AuthUsecase {
	return AuthUsecase{
		auth:  auth,
		users: users,
	}
}

func (u AuthUsecase) ValidateToken(token string) error {
	return u.auth.ValidateToken(token)
}

func (u AuthUsecase) SignUp(ctx context.Context, user domain.User) (string, error) {
	_, err := u.users.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	token, err := u.auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u AuthUsecase) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := u.users.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = user.CheckPassword(password)
	if err != nil {
		return "", err
	}
	token, err := u.auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}
