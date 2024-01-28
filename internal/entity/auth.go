package entity

import "errors"

type (
	User struct {
		UserID
		UserLogin
		UserPasswordHash
	}

	UserID           int32
	UserLogin        string
	UserPasswordHash string
	UserPassword     string

	Token string

	TokenData struct {
		UserID UserID
	}
)

var (
	ErrUserNotAvailable  = errors.New("login already taken")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
	ErrUnauthorized      = errors.New("authorization failed")
)
