package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nineee02/gotest/internal/dto"
	"github.com/nineee02/gotest/internal/service"
	"github.com/nineee02/gotest/pkg/app_context"
	"github.com/nineee02/gotest/pkg/validator"
)

type UserHandler interface {
	PostCreateUser(c echo.Context) error
	PostLogin(c echo.Context) error
}

type userHandler struct {
	userService service.UserService
	validator   *validator.Validator
}

func NewUserHandler(
	userService service.UserService,
	validator *validator.Validator,
) UserHandler {
	return &userHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *userHandler) PostCreateUser(c echo.Context) error {
	ctx := app_context.NewCtx(c).Build()
	
	var payload dto.UserRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status:  "error",
			Message: "Invalid request payload"})
	}

	if err := h.validator.Validate(payload); err != nil {
		if validatorErr, ok := validator.IsValidationErrors(err); ok {
			var fields []string
			for _, fieldErr := range validatorErr {
				fields = append(fields, fieldErr.Field())
			}
			return c.JSON(http.StatusUnprocessableEntity, dto.Response{
				Status:  "error",
				Message: "Validation failed",
				Data:    fields,
			})
		}
	}

	if err := h.userService.CreateUser(ctx, &payload); err != nil {
		if strings.Contains(err.Error(), "username") {
			return c.JSON(http.StatusConflict, dto.Response{
				Status:  "error",
				Message: "Username already exists"})
		}
		return c.JSON(http.StatusBadGateway, dto.Response{
			Status:  "error",
			Message: "Failed to create user",
			Errors:  err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Status:  "success",
		Message: "User created successfully"})
}

func (h *userHandler) PostLogin(c echo.Context) error {
	ctx := app_context.NewCtx(c).Build()

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	token, err := h.userService.Login(ctx, payload.Username, payload.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.Response{
			Status:  "error",
			Message: "Invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Status:  "success",
		Message: "Login successful",
		Data: map[string]string{
			"token": token,
		},
	})
}
