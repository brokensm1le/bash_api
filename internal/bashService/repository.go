package bashService

type Repository interface {
	GetCommand(cmdId int64) (*Command, error)
	GetList(params *GetListParams) ([]Command, error)
	CreateCommand(params *CreateCommandParams) (int64, error)
	DeleteCommand(id int64, userId int64) error
	DeleteCommandAdmin(id int64) error
}
