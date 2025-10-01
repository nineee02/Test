package repository

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/nineee02/gotest/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	var username string
	if username != "" {
		log.Errorf("Username already exists: %v", user.Username)
		return fmt.Errorf("username '%s' already exists", user.Username)
	}
	if err := r.db.WithContext(ctx).Model(&entity.User{}).
		Select("username").
		Where("username =? ", user.Username).
		Limit(1).
		Scan(&username).Error; err != nil {
		log.Errorf("Failed to check username: %v", err)
		return err
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Errorf("Failed to create user: %v", err)
		return err
	}

	return nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
