package domain

import "errors"

var (
	ErrTokenExpired   = errors.New("token is expired")
	ErrTokenBadClaims = errors.New("can not parse claims")
)
