package postgres

import (
	"context"
	"time"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
)

type user struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func userFromDomain(u domain.User) user {
	return user{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u user) userToDomain() domain.User {
	return domain.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

// CreateUser creates new user
func (c *PostgresClient) CreateUser(ctx context.Context, u domain.User) (*domain.User, error) {
	createUser := userFromDomain(u)
	stmt, err := c.db.PrepareNamedContext(ctx, `INSERT INTO users (username,email, password) 
	VALUES (:username, :email, :password) RETURNING *`)
	if err != nil {
		return nil, err
	}
	var newUser user
	err = stmt.GetContext(ctx, &newUser, createUser)
	domainUser := newUser.userToDomain()
	return &domainUser, err
}

func (c *PostgresClient) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u user
	err := c.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	domainUser := u.userToDomain()
	return &domainUser, nil
}
