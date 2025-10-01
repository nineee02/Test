package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nineee02/gotest/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/nineee02/gotest/pkg/config"
	"go.uber.org/dig"
)

const (
	defaultCtxTimeout     = 5 * time.Second
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultMaxHeaderBytes = 1 << 20 // 1 MB
)

type Server struct {
	cfg       *config.Configuration
	echo      *echo.Echo
	logger    logger.Logger
	container *dig.Container
}

func New(cfg *config.Configuration, logger logger.Logger) (*Server, error) {
	return &Server{
		cfg:       cfg,
		logger:    logger,
		echo:      echo.New(),
		container: dig.New(),
	}, nil
}

func (srv *Server) Start() error {
	if srv == nil {
		return fmt.Errorf("server instance is nil")
	}

	handlers, err := srv.initializeDependencies()
	if err != nil {
		return err
	}

	srv.registerRoutes(handlers, handlers.Middleware)
	go func() {
		httpServer := &http.Server{
			Addr: ":" + srv.cfg.Server.Port,
		}

		fmt.Println("✅ Server is listening on port:", srv.cfg.Server.Port)
		if err := srv.echo.StartServer(httpServer); err != nil {
			fmt.Println("❌ Server is shutting down:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return srv.Shutdown()
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.echo.Shutdown(ctx)
}
