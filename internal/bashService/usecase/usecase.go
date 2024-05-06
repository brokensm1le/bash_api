package usecase

import (
	"bash_api/internal/bashService"
	"bash_api/internal/cconstant"
)

type BashServiceUsecase struct {
	repo bashService.Repository
}

func NewBashServiceUsecase(repo bashService.Repository) bashService.Usecase {
	return &BashServiceUsecase{
		repo: repo,
	}
}

func (u *BashServiceUsecase) GetCommand(commandId int64) (*bashService.Command, error) {
	return u.repo.GetCommand(commandId)
}

func (u *BashServiceUsecase) GetList(params *bashService.GetListParams) ([]bashService.Command, error) {
	return u.repo.GetList(params)
}

func (u *BashServiceUsecase) CreateCommand(params *bashService.CreateCommandParams) (int64, error) {
	return u.repo.CreateCommand(params)
}

func (u *BashServiceUsecase) DeleteCommand(commandId int64, role int, userId int64) error {
	if role == cconstant.RoleAdmin {
		return u.repo.DeleteCommandAdmin(commandId)
	} else {
		return u.repo.DeleteCommand(commandId, userId)
	}
}
