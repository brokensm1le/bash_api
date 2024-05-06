package tokenManager

import "time"

type TokenManager interface {
	NewJWT(data *Data, ttl time.Duration) (string, error)
	Parse(accessToken string) (*Data, error)
	NewRefreshToken() (string, error)
}
