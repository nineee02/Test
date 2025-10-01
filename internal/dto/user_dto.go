package dto

import (
	"time"

	"github.com/nineee02/gotest/internal/entity"
)

type UserRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=255"`
	Email           string `json:"email" validate:"required,email,max=255"`
	Password        string `json:"password" validate:"required,min=3,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=3,max=255,eqfield=Password"`
}

func (u *UserRequest) ToUserRequestEntity() *entity.User {
	return &entity.User{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
