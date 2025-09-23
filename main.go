package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/route"
)

type Server struct {
	routers []route.Router
}

func (s *Server) start() {
	server := echo.New()

	for _, r := range s.routers {
		r.Configure(server)
	}

	server.Start(":8080")
}

func main() {
	ctx := context.Background()

	routers, err := route.Routers(ctx)
	if err != nil {
		log.Fatalf("Failed to get routers: %v", err)
	}

	s := Server{routers: routers}
	s.start()
}
