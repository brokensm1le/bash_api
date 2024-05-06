package usecase

import (
	"bash_api/internal/auth"
	"bash_api/internal/cconstant"
	"bash_api/pkg/customTime"
	"bash_api/pkg/hasher"
	"bash_api/pkg/tokenManager"
	"time"
)

type AuthUsecase struct {
	repo         auth.Repository
	hasher       hasher.PasswordHasher
	tokenManager tokenManager.TokenManager
}

func NewAuthUsecase(repo auth.Repository, hasher hasher.PasswordHasher, tokenManager tokenManager.TokenManager) auth.Usecase {
	return &AuthUsecase{
		repo:         repo,
		hasher:       hasher, // SHA256.NewSHA256Hasher(cconstant.Salt)
		tokenManager: tokenManager,
	}
}

func (u *AuthUsecase) SignUp(params *auth.SignUpParams) error {
	params.Password = u.hasher.Hash(params.Password)
	refreshToken, err := u.tokenManager.NewRefreshToken()
	if err != nil {
		return err
	}
	return u.repo.CreateUser(&auth.User{
		Name:            params.Name,
		Email:           params.Email,
		Password:        params.Password,
		RefreshToken:    refreshToken,
		RefreshTokenTTL: customTime.GetMoscowTime(),
	})
}

func (u *AuthUsecase) SignIn(params *auth.SignInParams) (*auth.TokensResponse, error) {
	var resp *auth.TokensResponse

	params.Password = u.hasher.Hash(params.Password)
	user, err := u.repo.GetUser(params)
	if err != nil {
		return nil, err
	}

	resp, err = u.createSession(user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *AuthUsecase) RefreshTokens(refreshToken string) (*auth.TokensResponse, error) {
	user, err := u.repo.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return u.createSession(user)
}

func (u *AuthUsecase) BeAdmin(id int64) error {
	return u.repo.BeAdmin(id)
}

// --------------------------------------------------------------------------------------------------------------------

func (u *AuthUsecase) createSession(user *auth.User) (*auth.TokensResponse, error) {
	var (
		res auth.TokensResponse
		err error
	)
	res.AccessToken, err = u.tokenManager.NewJWT(&tokenManager.Data{Id: user.Id, Role: user.Role}, cconstant.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	res.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = u.repo.SetRefreshToken(user.Id, res.RefreshToken, time.Now().Add(cconstant.RefreshTokenTTl).UTC())
	if err != nil {
		return nil, err
	}
	return &res, nil
}
