package bashService

type Usecase interface {
	GetCommand(commandId int64) (*Command, error)
	GetList(params *GetListParams) ([]Command, error)
	CreateCommand(params *CreateCommandParams) (int64, error)
	DeleteCommand(commandId int64, role int, personID int64) error
	RunCommand(commandId int64, personID int64) (int64, error)

	KillRun(personId int64, role int, runId int64) error
	GetPersonResult(params *GetListParams) ([]Result, error)
	GetRunResult(runId int64) (*Result, error)
}
