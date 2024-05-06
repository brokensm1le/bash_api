package tokenManager

import "github.com/dgrijalva/jwt-go"

type Data struct {
	Id   int64 `json:"id"`
	Role int   `json:"role"`
}

type CustomClaims struct {
	jwt.StandardClaims
	Id   int64 `json:"id"`
	Role int   `json:"role"`
}
