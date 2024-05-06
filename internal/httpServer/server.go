package httpServer

import (
	"bash_api/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	cfg   *config.Config
	fiber *fiber.App
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{DisableStartupMessage: true}),
		cfg:   cfg,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.fiber); err != nil {
		log.Fatalf("Cannot map handlers. Error: {%s}", err)
	}

	log.Printf("Start server on {host:port - %s:%s}", s.cfg.Server.Host, s.cfg.Server.Port)

	if err := s.fiber.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)); err != nil {
		log.Fatalf("Cannot listen. Error: {%s}", err)
	}

	return nil
}
