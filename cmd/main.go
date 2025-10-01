package main

import (
	"fmt"

	"github.com/nineee02/gotest/internal/server"
	"github.com/nineee02/gotest/pkg/config"
	"github.com/nineee02/gotest/pkg/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	appLogger := logger.NewZapLogger(cfg)

	srv, err := server.New(cfg, appLogger)
	if err != nil {
		panic(fmt.Errorf("failed to create server: %v", err))
	}

	if err := srv.Start(); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
