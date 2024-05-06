package hasher

type PasswordHasher interface {
	Hash(password string) string
}
