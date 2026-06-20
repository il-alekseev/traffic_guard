package keycloakclient

import "errors"

// Ошибки usecase
var (
	ErrTemporaryPassword      = errors.New("temporary password")
	ErrAccountNotSetUp        = errors.New("account not fully set up")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidUserCredentials = errors.New("invalid user credentials")
	ErrInvalidRole            = errors.New("invalid role")
)
