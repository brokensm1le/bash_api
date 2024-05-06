package bashService

import (
	"github.com/lib/pq"
	"time"
)

type Command struct {
	CmdId     int64          `json:"cmd_id" db:"cmd_id"`
	Cmd       string         `json:"cmd" db:"cmd"`
	CmdArgs   pq.StringArray `json:"cmd_args" db:"cmd_args"`
	AuthorId  int64          `json:"author_id" db:"author_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type CreateCommandParams struct {
	Cmd      string   `json:"cmd" db:"cmd"`
	CmdArgs  []string `json:"cmd_args" db:"cmd_args"`
	AuthorId int64    `json:"-" db:"author_id"`
}

type GetListParams struct {
	Limit    int64 `json:"limit"`
	Offset   int64 `json:"offset"`
	AuthorId int64 `json:"author_id"`
}
