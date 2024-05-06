package auth

import "time"

type Repository interface {
	CreateUser(user *User) error
	GetUser(params *SignInParams) (*User, error)
	SetRefreshToken(id int64, refresh string, refreshTTL time.Time) error
	GetByRefreshToken(refreshToken string) (*User, error)
	BeAdmin(id int64) error
}
