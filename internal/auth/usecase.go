package auth

type Usecase interface {
	SignIn(params *SignInParams) (*TokensResponse, error)
	SignUp(params *SignUpParams) error
	RefreshTokens(refreshToken string) (*TokensResponse, error)
	BeAdmin(id int64) error
}
