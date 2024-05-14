package test

import (
	"bash_api/internal/auth"
	http2 "bash_api/internal/auth/delivery/http"
	mockAuth "bash_api/internal/auth/mocks"
	"bash_api/internal/auth/usecase"
	"bash_api/internal/bashService"
	http3 "bash_api/internal/bashService/delivery/http"
	"bash_api/internal/bashService/fake_repository"
	usecase2 "bash_api/internal/bashService/usecase"
	"bash_api/internal/cconstant"
	"bash_api/pkg/hasher/SHA256"
	"bash_api/pkg/tokenManager/jwtTokenManager"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

type CreateCmdRes struct {
	Id int64 `json:"command_id"`
}

type CreateRunRes struct {
	Id int64 `json:"run_id"`
}

type CreateCommandParamsBad struct {
	BadField string `json:"badField"`
	AuthorId int64  `json:"-" db:"author_id"`
}

func Test(t *testing.T) {
	t.Run("testCGDCmd", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBashRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BashUC := usecase2.NewBashServiceUsecase(FakeBashRepo)
		handlerBash := http3.NewBashServiceHandler(BashUC, tokenManager)

		// Create test app
		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBash.UserIdentity())
		appWithToken.Get("/get_cmd/:cmd_id", handlerBash.GetCommand())
		appWithToken.Post("/create_cmd", handlerBash.CreateCommand())
		appWithToken.Delete("/delete/:cmd_id", handlerBash.DeleteCommand())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create bash
		r := bashService.CreateCommandParams{Cmd: "#!/bin/bash \n echo There were $# parameters passed.", CmdArgs: []string{"1", "2", "3", "4", "5"}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res CreateCmdRes
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, int64(1), res.Id)

		// TEST get bash
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_cmd/%d", res.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var resCmd bashService.Command
		err = json.Unmarshal(bodyBytes, &resCmd)
		require.NoError(t, err)
		require.Equal(t, r.AuthorId, resCmd.AuthorId)
		require.Equal(t, r.Cmd, resCmd.Cmd)
		require.Equal(t, len(r.CmdArgs), len(resCmd.CmdArgs))

		// TEST delete bash
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/delete/%d", res.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// TEST get bash
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_cmd/%d", res.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("testCmdError", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBashRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BashUC := usecase2.NewBashServiceUsecase(FakeBashRepo)
		handlerBash := http3.NewBashServiceHandler(BashUC, tokenManager)

		// Create test app
		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBash.UserIdentity())
		appWithToken.Get("/get_cmd/:cmd_id", handlerBash.GetCommand())
		appWithToken.Post("/create_cmd", handlerBash.CreateCommand())
		appWithToken.Delete("/delete/:cmd_id", handlerBash.DeleteCommand())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create bash (no accessToken)
		r := bashService.CreateCommandParams{Cmd: "#!/bin/bash \n echo There were $# parameters passed.", CmdArgs: []string{"1", "2", "3", "4", "5"}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// TEST create bash
		r = bashService.CreateCommandParams{Cmd: "#!/bin/bash \n echo There were $# parameters passed.", CmdArgs: []string{"1", "2", "3", "4", "5"}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res CreateCmdRes
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, int64(1), res.Id)

		// TEST get bash (bad request)
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_cmd/wer"), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// TEST get bash (not found)
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_cmd/%d", 99), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Test delete bash (not found)
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/delete/%d", 99), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("testCreateRunAndGetResult", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBashRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BashUC := usecase2.NewBashServiceUsecase(FakeBashRepo)
		handlerBash := http3.NewBashServiceHandler(BashUC, tokenManager)

		// Create test app

		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBash.UserIdentity())
		appWithToken.Post("/create_cmd", handlerBash.CreateCommand())
		appWithToken.Post("/run/:cmd_id", handlerBash.RunCommand())
		appWithToken.Get("/get_run/:run_id", handlerBash.GetRun())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create bash
		r := bashService.CreateCommandParams{Cmd: "#!/bin/bash \n echo There were $# parameters passed.", CmdArgs: []string{"1", "2", "3", "4", "5"}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res CreateCmdRes
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, int64(1), res.Id)

		// TEST create run
		req, _ = http.NewRequest("POST", fmt.Sprintf("/run/%d", res.Id), b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var resRun CreateRunRes
		err = json.Unmarshal(bodyBytes, &resRun)
		require.NoError(t, err)
		require.Equal(t, int64(1), resRun.Id)

		time.Sleep(time.Second * 1)

		// TEST get results run
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_run/%d", resRun.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		require.NoError(t, err)
		var runResult bashService.Result
		err = json.Unmarshal(bodyBytes, &runResult)
		require.NoError(t, err)
		require.Equal(t, "There were 5 parameters passed.\n", runResult.Results)

		// other test
		// TEST create bash
		r = bashService.CreateCommandParams{Cmd: "#!/bin/bash\nfor (( i=1; i <= 5; i++ ))\ndo\necho \"number is $i\"\ndone", CmdArgs: []string{}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, int64(2), res.Id)

		// TEST create run
		req, _ = http.NewRequest("POST", fmt.Sprintf("/run/%d", res.Id), b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = json.Unmarshal(bodyBytes, &resRun)
		require.NoError(t, err)
		require.Equal(t, int64(2), resRun.Id)

		time.Sleep(time.Second * 1)

		// TEST get results run
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_run/%d", resRun.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		require.NoError(t, err)
		err = json.Unmarshal(bodyBytes, &runResult)
		require.NoError(t, err)
		require.Equal(t, "number is 1\nnumber is 2\nnumber is 3\nnumber is 4\nnumber is 5\n", runResult.Results)
	})
	t.Run("testRunError", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		// ---- init
		hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
		mockAuthRepo := mockAuth.NewMockRepository(c)
		mockAuthRepo.EXPECT().GetUser(&auth.SignInParams{Email: "Sasha", Password: hasher.Hash("1234")}).Return(&auth.User{Id: 1, Role: 1}, nil).Times(1)
		mockAuthRepo.EXPECT().SetRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		FakeBashRepo := fake_repository.NewFakeRepository()

		tokenManager, _ := jwtTokenManager.NewManger(cconstant.SignedKey)
		AuthUC := usecase.NewAuthUsecase(mockAuthRepo, hasher, tokenManager)
		handlerAuth := http2.NewAuthHandler(AuthUC)

		BashUC := usecase2.NewBashServiceUsecase(FakeBashRepo)
		handlerBash := http3.NewBashServiceHandler(BashUC, tokenManager)

		// Create test app
		app := fiber.New()
		app.Post("/auth/signIn", handlerAuth.SignIn())
		appWithToken := app.Use("/", handlerBash.UserIdentity())
		appWithToken.Post("/create_cmd", handlerBash.CreateCommand())
		appWithToken.Post("/run/:cmd_id", handlerBash.RunCommand())
		appWithToken.Get("/get_run/:run_id", handlerBash.GetRun())

		// TEST get tokens
		m, b := map[string]string{"email": "Sasha", "password": "1234"}, new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(m)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/auth/signIn", b)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		tokens := auth.TokensResponse{}
		err = json.Unmarshal(bodyBytes, &tokens)
		require.NoError(t, err)

		// TEST create bash (bad cmd)
		r := bashService.CreateCommandParams{Cmd: "#!/bin/bash\nfor (( i=1; i <= 5; i++ ))\ndo\necho \"number is $i\"\n", CmdArgs: []string{}, AuthorId: 1}
		b = new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(r)
		require.NoError(t, err)
		req, _ = http.NewRequest("POST", "/create_cmd", b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var res CreateCmdRes
		err = json.Unmarshal(bodyBytes, &res)
		require.NoError(t, err)
		require.Equal(t, int64(1), res.Id)

		// TEST create run (bad cmd)
		req, _ = http.NewRequest("POST", fmt.Sprintf("/run/%d", res.Id), b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		var resRun CreateRunRes
		err = json.Unmarshal(bodyBytes, &resRun)
		require.NoError(t, err)
		require.Equal(t, int64(1), resRun.Id)

		time.Sleep(time.Second * 1)

		// TEST get results run (bad cmd)
		req, _ = http.NewRequest("GET", fmt.Sprintf("/get_run/%d", resRun.Id), nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, err = io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		require.NoError(t, err)
		var runResult bashService.Result
		err = json.Unmarshal(bodyBytes, &runResult)
		require.NoError(t, err)
		require.Equal(t, "exit status 2", runResult.Results)

		// TEST create run (bad id)
		req, _ = http.NewRequest("POST", fmt.Sprintf("/run/dck"), b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// TEST create run (dont found id)
		req, _ = http.NewRequest("POST", fmt.Sprintf("/run/%d", 10), b)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		resp, err = app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
