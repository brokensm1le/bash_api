package auth

import "time"

type User struct {
	Id              int64     `json:"-" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"password" db:"password"`
	Role            int       `json:"-" db:"role"`
	RefreshToken    string    `json:"-" db:"refresh_token"`
	RefreshTokenTTL time.Time `json:"-" db:"refresh_token_ttl"`
}

type SignUpParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
