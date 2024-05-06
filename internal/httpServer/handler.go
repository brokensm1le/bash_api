package httpServer

import (
	"bash_api/internal/auth/delivery/http"
	authRepository "bash_api/internal/auth/repository"
	"bash_api/internal/auth/usecase"
	http2 "bash_api/internal/bashService/delivery/http"
	bashServiceRepository "bash_api/internal/bashService/repository"
	usecase2 "bash_api/internal/bashService/usecase"
	"bash_api/internal/cconstant"
	"bash_api/pkg/hasher/SHA256"
	"bash_api/pkg/storagePostgres"
	"bash_api/pkg/tokenManager/jwtTokenManager"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (s *Server) MapHandlers(app *fiber.App) error {
	db, err := storagePostgres.InitPsqlDB(s.cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = storagePostgres.CreateTable(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	hasher := SHA256.NewSHA256Hasher(cconstant.Salt)
	manager, err := jwtTokenManager.NewManger(cconstant.SignedKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	authRepo := authRepository.NewPostgresRepository(db)
	bashRepo := bashServiceRepository.NewPostgresRepository(db)

	authUC := usecase.NewAuthUsecase(authRepo, hasher, manager)
	bashUC := usecase2.NewBashServiceUsecase(bashRepo)

	authR := http.NewAuthHandler(authUC)
	bashR := http2.NewBashServiceHandler(bashUC, manager)

	http.MapRoutes(app, authR)
	http2.MapRoutes(app, bashR)

	return nil
}
