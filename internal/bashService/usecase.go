package bashService

type Usecase interface {
	GetCommand(commandId int64) (*Command, error)
	GetList(params *GetListParams) ([]Command, error)
	CreateCommand(params *CreateCommandParams) (int64, error)
	DeleteCommand(commandId int64, role int, userId int64) error
}
