package service

import (
	"errors"
	"fmt"
	"simple-emoney/config"
	"simple-emoney/internal/app/repository"
	"simple-emoney/internal/model"
	"simple-emoney/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req *model.RegisterRequest) (*model.User, error)
	LoginUser(req *model.LoginRequest) (*model.LoginResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewAuthService(userRepo repository.UserRepository, redisRepo repository.RedisRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *authService) RegisterUser(req *model.RegisterRequest) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Balance:      0,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	err = s.redisRepo.SetUserCache(user, 1*time.Hour)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Username, s.cfg.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	err = s.redisRepo.SetAuthToken(token, user.ID.String(), 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Token:    token,
	}, nil
}
