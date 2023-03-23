package service

import (
	"time"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/dgrijalva/jwt-go"
)

const expireAt = time.Hour

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type AuthService struct {
	jwtKey   []byte
	expireAt time.Duration
}

func NewAuthService(jwtKey []byte) *AuthService {
	return &AuthService{
		jwtKey:   jwtKey,
		expireAt: expireAt,
	}
}

// GenerateJWT generate JWT for email and username
func (a *AuthService) GenerateJWT(email, username string) (string, error) {
	expirationTime := time.Now().Add(a.expireAt)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtKey)
}

// ValidateToken validates token
func (a *AuthService) ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return a.jwtKey, nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return domain.ErrTokenBadClaims
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return domain.ErrTokenExpired
	}
	return nil
}
