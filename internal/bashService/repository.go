package bashService

type Repository interface {
	GetCommand(cmdId int64) (*Command, error)
	GetList(params *GetListParams) ([]Command, error)
	CreateCommand(params *CreateCommandParams) (int64, error)
	DeleteCommand(id int64, userId int64) error
	DeleteCommandAdmin(id int64) error

	CreateRun(params *CreateRunParams) (int64, error)
	ChangeRunStatus(params *ChngRunStatusParams) error
	GetRun(runId int64) (*Result, error)
	GetPersonRun(params *GetListParams) ([]Result, error)
	GetAuthorIdByRunId(runId int64) (int64, error)
}
