package middleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nineee02/gotest/pkg/constant"
	"go.uber.org/zap"
)

func (mw *MiddlewareManager) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := c.Request().Header.Get(echo.HeaderXRequestID)
		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Response().Header().Set(echo.HeaderXRequestID, requestId)
		c.Set(echo.HeaderXRequestID, requestId)

		// ✅ แก้ตรงนี้: set ลง context ด้วย
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, constant.RequestIdKey, requestId)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (mw *MiddlewareManager) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		start := time.Now()

		requestId := req.Header.Get(echo.HeaderXRequestID)
		if requestId == "" {
			requestId = c.Response().Header().Get(echo.HeaderXRequestID)
		}

		mw.logger.
			WithFields(
				zap.String(constant.RequestIdKey, requestId),
				zap.String("method", req.Method),
				zap.String("ip", c.RealIP()),
				zap.String("path", req.URL.Path),
			).
			Info("API Call Request")

		err := next(c)
		res := c.Response()

		if err != nil {
			mw.logger.
				WithFields(
					zap.String(constant.RequestIdKey, requestId),
					zap.Int("status", res.Status),
					zap.Int64("latency_ms", time.Since(start).Milliseconds()),
				).
				Error("API Call Response")
		} else {
			mw.logger.
				WithFields(
					zap.String(constant.RequestIdKey, requestId),
					zap.Int("status", res.Status),
					zap.Int64("latency_ms", time.Since(start).Milliseconds()),
				).
				Info("API Call Response")
		}

		return err
	}
}
