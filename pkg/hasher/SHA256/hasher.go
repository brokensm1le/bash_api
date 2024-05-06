package SHA256

import (
	"bash_api/pkg/hasher"
	"crypto/sha256"
	"fmt"
)

type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) hasher.PasswordHasher {
	return &SHA256Hasher{salt: salt}
}

func (h *SHA256Hasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
