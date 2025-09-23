package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lamkn06/user-app-golang.git/internal/route"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"

	"github.com/labstack/echo/v4"
)

var runtimeConfig runtime.ServerConfig

type Server struct {
	config  runtime.ServerConfig
	routers []route.Router
}

func (s *Server) start() {
	server := echo.New()

	for _, r := range s.routers {
		r.Configure(server)
	}

	log.Println("Server started on port========", s.config)
	channel := make(chan error)
	go func() {
		channel <- server.Start(":" + s.config.Port)
	}()
	log.Println("Server started on port", s.config.Port)
	select {
	case sig := <-shutdownSignals():
		log.Println("Shutting down server...")
		log.Printf("Received signal: %v", sig)
	case err := <-channel:
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	ctx := context.Background()

	runtime.LoadConfigs([]any{&runtimeConfig})

	fmt.Println("runtimeConfig========", runtimeConfig.Port)
	routers, err := route.Routers(ctx, runtimeConfig)
	if err != nil {
		log.Fatalf("Failed to get routers: %v", err)
	}

	s := Server{routers: routers, config: runtimeConfig}
	s.start()
}

func shutdownSignals() (signals <-chan os.Signal) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGABRT, syscall.SIGTERM)
	return channel
}
