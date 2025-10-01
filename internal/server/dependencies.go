package server

import (
	"fmt"

	"github.com/nineee02/gotest/pkg/logger"

	"github.com/nineee02/gotest/internal/handler"
	"github.com/nineee02/gotest/internal/repository"
	"github.com/nineee02/gotest/internal/service"
	"github.com/nineee02/gotest/pkg/config"
	"github.com/nineee02/gotest/pkg/database"
	"github.com/nineee02/gotest/pkg/middleware"
	"github.com/nineee02/gotest/pkg/util"
	"github.com/nineee02/gotest/pkg/validator"
	"gorm.io/gorm"
)

func (s *Server) initializeDependencies() (*HandlerGroup, error) {

	s.container.Provide(func() *config.Configuration {
		cfg, err := config.New()
		if err != nil {
			panic(err)
		}
		return cfg
	})

	s.container.Provide(func(cfg *config.Configuration) (*gorm.DB, error) {
		srv, err := database.NewMySQLDB(cfg)
		if err != nil {
			fmt.Println("failed to connect to mysql: ", err)
			return nil, err
		}
		fmt.Println("✅ connected to mysql successfully")
		return srv, nil
	})

	s.container.Provide(func() util.AESUtil {
		return &util.AESUtilImpl{}
	})

	s.container.Provide(handler.NewHealthHandler)
	s.container.Provide(validator.NewValidator)

	s.container.Provide(func(cfg *config.Configuration, log logger.Logger) *middleware.MiddlewareManager {
		return middleware.NewMiddlewareManager(cfg, log)
	})

	s.container.Provide(repository.NewUserRepository)

	s.container.Provide(service.NewUserService)

	s.container.Provide(handler.NewUserHandler)

	var handlers HandlerGroup
	if err := s.container.Invoke(func(
		user handler.UserHandler,
		health handler.HealthHandler,
		validator *validator.Validator,
	) {
		handlers = HandlerGroup{
			User:      user,
			Health:    health,
			Validator: validator,
		}
	}); err != nil {
		return nil, err
	}

	// s.applyMiddleware(handlers.Middleware)

	s.echo.Validator = handlers.Validator

	return &handlers, nil
}

func (s *Server) applyMiddleware(mw *middleware.MiddlewareManager) {
	s.echo.Use(
		mw.RequestID,
		mw.AcceptLanguage,
		mw.Logger,
		mw.BodyDump(),
		mw.CORSWithConfig(s.cfg),
		mw.TimeoutWithConfig(s.cfg),
		mw.Recover(),
	)
}
