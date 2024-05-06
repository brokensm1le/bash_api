package repository

import (
	"bash_api/internal/bashService"
	"bash_api/internal/cconstant"
	"bash_api/pkg/customTime"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) bashService.Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) GetCommand(cmdId int64) (*bashService.Command, error) {
	var (
		data  bashService.Command
		query = `
		SELECT *
		FROM %[1]s 
		WHERE cmdId = $1;
		`

		values = []any{cmdId}
	)

	query = fmt.Sprintf(query, cconstant.CommandsDB)

	if err := p.db.Get(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetList(params *bashService.GetListParams) ([]bashService.Command, error) {
	var (
		data  []bashService.Command
		query = `
		SELECT *
		FROM %[1]s
		`

		values = []any{params.Limit, params.Offset}
	)

	query = fmt.Sprintf(query, cconstant.CommandsDB)

	if params.AuthorId != -1 {
		query += "WHERE author_id = $3"
		values = append(values, params.AuthorId)
	}
	query += "LIMIT $1 OFFSET $2;"

	if err := p.db.Select(&data, query, values...); err != nil {
		return data, err
	}

	return data, nil
}

func (p *postgresRepository) CreateCommand(params *bashService.CreateCommandParams) (int64, error) {
	var (
		query = `INSERT INTO %[1]s (cmd, cmd_args, author_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING cmd_id;`

		timeNow time.Time = customTime.GetMoscowTime()
		values            = []any{params.Cmd, pq.Array(params.CmdArgs), params.AuthorId, timeNow}
		id      int64
	)

	query = fmt.Sprintf(query, cconstant.CommandsDB)

	if err := p.db.Get(&id, query, values...); err != nil {
		return id, err
	}

	return id, nil
}

func (p *postgresRepository) DeleteCommand(id int64, userId int64) error {
	var (
		query = `
		DELETE FROM %[1]s 
		WHERE cmd_id = $1 && authorId = $2
		`

		values = []any{id, userId}
	)

	query = fmt.Sprintf(query, cconstant.CommandsDB)

	res, err := p.db.Exec(query, values...)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("this is not your script or it has already been deleted")
	}

	return nil
}

func (p *postgresRepository) DeleteCommandAdmin(id int64) error {
	var (
		query = `
		DELETE FROM %[1]s 
		WHERE cmd_id = $1
		`

		values = []any{id}
	)

	query = fmt.Sprintf(query, cconstant.CommandsDB)

	res, err := p.db.Exec(query, values...)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("this is script has already been deleted")
	}

	return nil
}