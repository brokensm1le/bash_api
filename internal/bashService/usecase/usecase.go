package usecase

import (
	"bash_api/internal/bashService"
	"bash_api/internal/cconstant"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
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

func (u *BashServiceUsecase) RunCommand(commandId int64, userId int64) (int64, error) {
	command, err := u.repo.GetCommand(commandId)
	if err != nil {
		return 0, err
	}

	//u.repo.CreateRun()
	nameFile := fmt.Sprintf("%d_%d.sh", command.CmdId, time.Now().Unix())

	f, err := os.Create(nameFile)
	if err != nil {
		log.Println(err)
	}
	f.WriteString(command.Cmd)
	f.Close()

	stdout, err := exec.Command("chmod", "+x", nameFile).Output()
	log.Println(string(stdout))
	stdout, err = exec.Command("cat", nameFile).Output()
	log.Println(string(stdout))
	stdout, err = exec.Command("ls", "-la").Output()
	log.Println(string(stdout))

	command.CmdArgs = append([]string{"./" + nameFile}, command.CmdArgs...)

	log.Println(command.CmdArgs)

	stdout, err = exec.Command("/bin/sh", command.CmdArgs...).Output()
	log.Println(string(stdout))

	return 0, nil
}

func (u *BashServiceUsecase) GetRunResult(runId int64) {
}

func (u *BashServiceUsecase) GetPersonResult() {
}
