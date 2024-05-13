package usecase

import (
	"bash_api/internal/bashService"
	"bash_api/internal/cconstant"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type BashServiceUsecase struct {
	repo      bashService.Repository
	cancelMap sync.Map
}

func NewBashServiceUsecase(repo bashService.Repository) bashService.Usecase {
	return &BashServiceUsecase{
		repo:      repo,
		cancelMap: sync.Map{},
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

func (u *BashServiceUsecase) DeleteCommand(commandId int64, role int, personID int64) error {
	if role == cconstant.RoleAdmin {
		return u.repo.DeleteCommandAdmin(commandId)
	} else {
		return u.repo.DeleteCommand(commandId, personID)
	}
}

func (u *BashServiceUsecase) RunCommand(commandId int64, personID int64) (int64, error) {
	command, err := u.repo.GetCommand(commandId)
	if err != nil {
		return 0, err
	}

	runId, err := u.repo.CreateRun(&bashService.CreateRunParams{CmdId: commandId, AuthorId: personID})
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	u.cancelMap.Store(runId, cancel)

	go func() {
		nameFile := fmt.Sprintf("%d.sh", runId)

		f, err := os.Create(nameFile)
		if err != nil {
			log.Println("Error create file:", err)
			err := u.repo.ChangeRunStatus(&bashService.ChngRunStatusParams{RunId: runId, StatusId: cconstant.FailedStatus, Result: err.Error()})
			if err != nil {
				log.Println("Error in run command -> ChangeRunStatus:", err)
			}
			return
		}
		_, err = f.WriteString(command.Cmd)
		if err != nil {
			log.Println("Error write file:", err)
			err := u.repo.ChangeRunStatus(&bashService.ChngRunStatusParams{RunId: runId, StatusId: cconstant.FailedStatus, Result: err.Error()})
			if err != nil {
				log.Println("Error in run command -> ChangeRunStatus:", err)
			}
			return
		}
		f.Close()

		_, err = exec.Command("chmod", "+x", nameFile).Output()
		if err != nil {
			log.Println("Error in chmod file:", err)
			err := u.repo.ChangeRunStatus(&bashService.ChngRunStatusParams{RunId: runId, StatusId: cconstant.FailedStatus, Result: err.Error()})
			if err != nil {
				log.Println("Error in run command -> ChangeRunStatus:", err)
			}
			return
		}

		command.CmdArgs = append([]string{"./" + nameFile}, command.CmdArgs...)

		stdout, err := exec.CommandContext(ctx, "/bin/sh", command.CmdArgs...).Output()
		if err != nil {
			log.Println("Error in run command:", err)
			err := u.repo.ChangeRunStatus(&bashService.ChngRunStatusParams{RunId: runId, StatusId: cconstant.FailedStatus, Result: err.Error()})
			if err != nil {
				log.Println("Error in run command -> ChangeRunStatus:", err)
			}
			return
		}

		err = os.Remove(nameFile)
		if err != nil {
			log.Println("Error in delete file:", err)
		}
		err = u.repo.ChangeRunStatus(&bashService.ChngRunStatusParams{RunId: runId, StatusId: cconstant.SuccessStatus, Result: string(stdout)})
		if err != nil {
			log.Println("Error in ChangeRunStatus:", err)
		}
	}()

	return runId, nil
}

func (u *BashServiceUsecase) GetRunResult(runId int64) (*bashService.Result, error) {
	return u.repo.GetRun(runId)
}

func (u *BashServiceUsecase) GetPersonResult(params *bashService.GetListParams) ([]bashService.Result, error) {
	return u.repo.GetPersonRun(params)
}

func (u *BashServiceUsecase) KillRun(personID int64, runId int64) error {
	// check_person
	value, ok := u.cancelMap.Load(runId)
	if !ok {
		return fmt.Errorf("process not found! =(")
	}
	cancelF := value.(context.CancelFunc)
	cancelF()
	u.cancelMap.Delete(runId)
	return nil
}
