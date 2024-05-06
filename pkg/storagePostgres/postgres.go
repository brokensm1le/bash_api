package storagePostgres

import (
	"bash_api/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

// ------------------------------------------------------------------------------------------------------------------------------

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SSLMode)

	return sqlx.Connect(c.Postgres.PgDriver, connectionUrl)
}

func CreateTable(db *sqlx.DB) error {
	var (
		query = `
		CREATE TABLE IF NOT EXISTS "cmds"
		(
			cmd_id       	bigserial    not null unique,
			cmd			   	text       	 not null,
			cmd_args   		text[]	 	 ,
			author_id		bigint       not null,
			created_at      timestamp	 not null
		);
		CREATE TABLE IF NOT EXISTS "results"
		(
		    run_id			bigserial    not null unique,
			cmd_id       	bigint    	 not null,
			author_id		bigint       not null,
			results			text		 ,
			created_at      timestamp	 not null
		);
		CREATE TABLE IF NOT EXISTS "auth"
		(
			id		  	    	bigserial    not null unique,
			name				text		 not null,
			email		   		varchar(255) not null unique,
			password   			text		 not null,
			role      			smallint	 not null default 0,
			refresh_token    	text		 not null unique,
			refresh_token_ttl   timestamp	 not null
		);
		`
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
