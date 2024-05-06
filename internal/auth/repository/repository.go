package repository

import (
	"bash_api/internal/auth"
	"bash_api/internal/cconstant"
	"bash_api/pkg/customTime"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) auth.Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) CreateUser(user *auth.User) error {
	var (
		query = `
		INSERT INTO %[1]s (name, email, password, refresh_token, refresh_token_ttl)
		VALUES ($1, $2, $3, $4, $5)`

		values = []any{user.Name, user.Email, user.Password, user.RefreshToken, user.RefreshTokenTTL}
	)

	query = fmt.Sprintf(query, cconstant.AuthDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetUser(params *auth.SignInParams) (*auth.User, error) {
	var (
		data  []auth.User
		query = `
		SELECT *
		FROM %[1]s
		WHERE email=$1 AND password=$2
		`

		values = []any{params.Email, params.Password}
	)

	query = fmt.Sprintf(query, cconstant.AuthDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &auth.User{}, err
	}

	if len(data) == 0 {
		return &auth.User{}, fmt.Errorf("uncorrect login or password")
	}

	return &data[0], nil
}

func (p *postgresRepository) SetRefreshToken(id int64, refresh string, refreshTTL time.Time) error {
	var (
		query string = `
		UPDATE %[1]s SET (refresh_token, refresh_token_ttl) =
			($1, $2)
		WHERE id = $3;
		`
		values []any = []any{refresh, refreshTTL, id}
	)

	// -----------------------------------------------------------------------------------------------------------------------------

	query = fmt.Sprintf(query, cconstant.AuthDB)

	// -----------------------------------------------------------------------------------------------------------------------------

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetByRefreshToken(refreshToken string) (*auth.User, error) {
	var (
		data  []auth.User
		query = `
		SELECT *
		FROM %[1]s
		WHERE refresh_token = $1 AND refresh_token_ttl > $2
		`

		values = []any{refreshToken, customTime.GetMoscowTime()}
	)

	query = fmt.Sprintf(query, cconstant.AuthDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &auth.User{}, err
	}

	if len(data) == 0 {
		return &auth.User{}, fmt.Errorf("complete the sign in again.")
	}

	return &data[0], nil
}

func (p *postgresRepository) BeAdmin(id int64) error {
	var (
		query string = `
		UPDATE %[1]s SET (role) =
			($1)
		WHERE id = $3;
		`
		values []any = []any{cconstant.RoleAdmin, id}
	)

	// -----------------------------------------------------------------------------------------------------------------------------

	query = fmt.Sprintf(query, cconstant.AuthDB)

	// -----------------------------------------------------------------------------------------------------------------------------

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}
