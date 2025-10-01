package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/nineee02/gotest/internal/dto"
	"github.com/nineee02/gotest/internal/repository"
	"github.com/nineee02/gotest/pkg/config"
	"github.com/nineee02/gotest/pkg/util"
)

type UserService interface {
	CreateUser(ctx context.Context, user *dto.UserRequest) error
	Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
	util           util.AESUtil
	cfg            *config.Configuration
}

func NewUserService(
	userRepository repository.UserRepository,
	util util.AESUtil,
	cfg *config.Configuration,
) UserService {
	return &userService{
		userRepository: userRepository,
		util:           util,
		cfg:            cfg,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *dto.UserRequest) error {

	if user.Password != user.ConfirmPassword {
		return fmt.Errorf("password and confirm password do not match")
	}

	encrypted, err := s.util.AES256Encrypt([]byte(user.Password), []byte(s.cfg.API.AesKey))
	if err != nil {
		return err
	}

	userReq := user.ToUserRequestEntity()
	userReq.Password = encrypted
	userReq.UserID = uuid.New()

	if err := s.userRepository.CreateUser(ctx, userReq); err != nil {
		log.Errorf("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Decrypt and compare passwords
	decryptedPassword, err := s.util.AES256Decrypt(user.Password, []byte(s.cfg.API.AesKey))
	if err != nil || string(decryptedPassword) != password {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := util.GenerateJWT(user.UserID.String(), s.cfg.API.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
