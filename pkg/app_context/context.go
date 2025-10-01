package app_context

import (
	"context"

	"github.com/labstack/echo/v4"
)

const (
	// RequestIDKey is the key for the request ID in the context
	RequestIDKey = "request_id"
	// UserIDKey is the key for the user ID in the context
	UserIDKey = "user_id"
	// RoleKey is the key for the role in the context
	RoleKey = "role"
	// EmailKey is the key for the email in the context
	EmailKey = "email"
	// FullNameKey is the key for the full name in the context
	FullNameKey = "full_name"
	// UsernameKey is the key for the username in the context
	UsernameKey = "username"
)

type AppContext struct {
	context.Context
	EchoCtx   echo.Context
	RequestID string
	UserID    string
	Role      string
	Email     string
	FullName  string
	Username  string
}

type AppContextBuilder struct {
	appContext *AppContext
}

func NewAppContext(c echo.Context) *AppContext {
	return &AppContext{Context: c.Request().Context(), EchoCtx: c}
}

func NewCtx(c echo.Context) *AppContextBuilder {
	return &AppContextBuilder{appContext: NewAppContext(c)}
}

func WithCtx(ctx context.Context) *AppContext {
	return &AppContext{
		Context:   ctx,
		RequestID: getStringFromCtx(ctx, RequestIDKey),
		UserID:    getStringFromCtx(ctx, UserIDKey),
		Role:      getStringFromCtx(ctx, RoleKey),
		Email:     getStringFromCtx(ctx, EmailKey),
		FullName:  getStringFromCtx(ctx, FullNameKey),
		Username:  getStringFromCtx(ctx, UsernameKey),
	}
}

func getStringFromCtx(ctx context.Context, key string) string {

	if val := ctx.Value(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (b *AppContextBuilder) RequestID() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(RequestIDKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, RequestIDKey, value.(string))
		b.appContext.RequestID = value.(string)
	}
	return b
}

func (b *AppContextBuilder) UserID() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(UserIDKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, UserIDKey, value.(string))
		b.appContext.UserID = value.(string)
	}
	return b
}

func (b *AppContextBuilder) Role() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(RoleKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, RoleKey, value.(string))
		b.appContext.Role = value.(string)
	}
	return b
}

func (b *AppContextBuilder) Email() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(EmailKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, EmailKey, value.(string))
		b.appContext.Email = value.(string)
	}
	return b
}

func (b *AppContextBuilder) FullName() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(FullNameKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, FullNameKey, value.(string))
		b.appContext.FullName = value.(string)
	}
	return b
}

func (b *AppContextBuilder) Username() *AppContextBuilder {
	value := b.appContext.EchoCtx.Get(UsernameKey)
	if value != nil {
		b.appContext.Context = context.WithValue(b.appContext.Context, UsernameKey, value.(string))
		b.appContext.Username = value.(string)
	}
	return b
}

func (b *AppContextBuilder) Build() *AppContext {
	return b.appContext
}

// Getter Methods
func (ac *AppContext) GetUserID() string {
	return ac.UserID
}

func (ac *AppContext) GetRole() string {
	return ac.Role
}

func (ac *AppContext) GetEmail() string {
	return ac.Email
}

func (ac *AppContext) GetFullName() string {
	return ac.FullName
}

func (ac *AppContext) GetUsername() string {
	return ac.Username
}
