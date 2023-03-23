package domain

import "golang.org/x/crypto/bcrypt"

// User is user of service
type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}

// HashPassword hashs password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword check password
func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
