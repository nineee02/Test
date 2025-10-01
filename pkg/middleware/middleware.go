package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nineee02/gotest/pkg/config"
	"github.com/nineee02/gotest/pkg/logger"
	"go.uber.org/zap"
)

type MiddlewareManager struct {
	cfg *config.Configuration

	logger logger.Logger
}

func NewMiddlewareManager(cfg *config.Configuration, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		logger: logger,
	}
}

func (mw *MiddlewareManager) Recover() echo.MiddlewareFunc {
	return echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			requestId := c.Get(echo.HeaderXRequestID).(string)
			mw.logger.With(
				zap.String("request_id", requestId),
				zap.ByteString("stack", stack),
				zap.Error(err),
			).Error(c.Request().Context(), "Panic Recover")
			return err
		},
	})
}

func (mw *MiddlewareManager) BodyDump() echo.MiddlewareFunc {
	return echoMiddleware.BodyDump(func(c echo.Context, reqBody, respBody []byte) {
		requestId := c.Get(echo.HeaderXRequestID).(string)
		requestPath := c.Request().RequestURI
		if strings.HasPrefix(requestPath, "/uploads/") {
			mw.logger.With(
				zap.String("request_id", requestId),
				zap.String("file_accessed", requestPath),
				zap.String("response_status", fmt.Sprintf("%d", c.Response().Status)),
			).Info(c.Request().Context(), "Static File Access")
			return
		}

		contentType := c.Request().Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {

			if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
				mw.logger.With(
					zap.String("request_id", requestId),
					zap.String("request_body", "[FAILED TO PARSE FORM DATA]"),
					zap.String("response_body", string(respBody)),
				).Warn(c.Request().Context(), "Body")
				return
			}

			formParams := make(map[string]string)
			for key, values := range c.Request().MultipartForm.Value {
				formParams[key] = values[0] // ใช้ค่าแรกของแต่ละคีย์
			}

			for key := range c.Request().MultipartForm.File {
				formParams[key] = "[FILE UPLOADED]"
			}

			mw.logger.With(
				zap.String("request_id", requestId),
				zap.Any("request_body", formParams),
				zap.String("response_body", string(respBody)),
			).Info(c.Request().Context(), "Body")
			return
		}

		mw.logger.With(
			zap.String("request_id", requestId),
			zap.String("request_body", string(reqBody)),
			zap.String("response_body", string(respBody)),
		).Info(c.Request().Context(), "Body")
	})
}

func (mw *MiddlewareManager) CORSWithConfig(cfg *config.Configuration) echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		// AllowOrigins:     cfg.Server.AllowOrigins,
		AllowHeaders:     cfg.Server.AllowHeaders,
		AllowMethods:     cfg.Server.AllowMethods,
		AllowCredentials: true,
	})
}

func (mw *MiddlewareManager) TimeoutWithConfig(cfg *config.Configuration) echo.MiddlewareFunc {
	return echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout: 5 * time.Second,
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/ws/")
		},
	})
}
