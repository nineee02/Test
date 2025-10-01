package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/nineee02/gotest/pkg/constant"
)

func (mw *MiddlewareManager) AcceptLanguage(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := c.Request().Header.Get("Accept-Language")
		if lang == "" {
			lang = "th"
		}
		c.Set(constant.AcceptLanguage, lang)
		return next(c)
	}
}
