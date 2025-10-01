package server

import (
	"context"

	"github.com/nineee02/gotest/internal/handler"
	"github.com/nineee02/gotest/pkg/middleware"
	"github.com/nineee02/gotest/pkg/validator"
	"go.uber.org/zap"
)

type HandlerGroup struct {
	Middleware *middleware.MiddlewareManager
	Validator  *validator.Validator
	Health     handler.HealthHandler
	User       handler.UserHandler
}

func (s *Server) registerRoutes(h *HandlerGroup, mw *middleware.MiddlewareManager) {

	s.echo.GET("/health", h.Health.Check)
	v1 := s.echo.Group("/api/v1")

	v1.POST("/users/register", h.User.PostCreateUser)
	v1.POST("/users/login", h.User.PostLogin)

	for _, route := range s.echo.Routes() {
		s.logger.Debug(context.TODO(), "Register route", zap.String("method", route.Method), zap.String("path", route.Path))
	}
}
