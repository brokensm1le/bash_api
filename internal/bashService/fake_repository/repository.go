package fake_repository

import (
	"bash_api/internal/bashService"
	"bash_api/internal/cconstant"
	"fmt"
)

type fakeRepository struct {
	commandId int64
	resultId  int64
	CommandDB map[int64]bashService.Command
	ResultDB  map[int64]bashService.Result
}

func NewFakeRepository() bashService.Repository {
	return &fakeRepository{
		commandId: 0,
		resultId:  0,
		CommandDB: make(map[int64]bashService.Command),
		ResultDB:  make(map[int64]bashService.Result),
	}
}

func (f *fakeRepository) GetCommand(cmdId int64) (*bashService.Command, error) {
	res, ok := f.CommandDB[cmdId]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &res, nil
}

func (f *fakeRepository) GetList(params *bashService.GetListParams) ([]bashService.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakeRepository) CreateCommand(params *bashService.CreateCommandParams) (int64, error) {
	f.commandId++
	f.CommandDB[f.commandId] = bashService.Command{
		AuthorId: params.AuthorId,
		Cmd:      params.Cmd,
		CmdArgs:  params.CmdArgs,
		CmdId:    f.commandId,
	}
	return f.commandId, nil
}

func (f fakeRepository) DeleteCommand(id int64, userId int64) error {
	res, ok := f.CommandDB[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	if res.AuthorId != userId {
		return fmt.Errorf("you aren't creator")
	}
	delete(f.CommandDB, id)
	return nil
}

func (f *fakeRepository) DeleteCommandAdmin(id int64) error {
	delete(f.CommandDB, id)
	return nil
}

func (f *fakeRepository) CreateRun(params *bashService.CreateRunParams) (int64, error) {
	f.resultId++
	f.ResultDB[f.resultId] = bashService.Result{
		RunId:    f.resultId,
		CmdId:    params.CmdId,
		AuthorId: params.AuthorId,
		StatusId: cconstant.RunningStatus,
	}
	return f.resultId, nil
}

func (f *fakeRepository) ChangeRunStatus(params *bashService.ChngRunStatusParams) error {
	res, ok := f.ResultDB[params.RunId]
	if !ok {
		return fmt.Errorf("not found")
	}

	res.StatusId = int64(params.StatusId)
	res.Results += params.Result
	f.ResultDB[params.RunId] = res

	return nil
}

func (f *fakeRepository) GetRun(runId int64) (*bashService.Result, error) {
	res, ok := f.ResultDB[runId]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &res, nil
}

func (f *fakeRepository) GetPersonRun(params *bashService.GetListParams) ([]bashService.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakeRepository) GetAuthorIdByRunId(runId int64) (int64, error) {
	res, ok := f.ResultDB[runId]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return res.AuthorId, nil
}
