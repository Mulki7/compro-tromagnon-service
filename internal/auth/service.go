package auth

import (
	"errors"
	"time"

	"compro-tromagnon-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(req LoginRequest) (*LoginResponse, error)
	GetProfile(adminID uuid.UUID) (*AdminData, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	admin, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	expiry := time.Now().Add(time.Hour * time.Duration(s.cfg.JWTExpiryHours))
	claims := jwt.MapClaims{
		"admin_id": admin.ID.String(),
		"email":    admin.Email,
		"name":     admin.Name,
		"exp":      expiry.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}

	return &LoginResponse{
		Token: tokenString,
		Admin: AdminData{
			ID:    admin.ID,
			Email: admin.Email,
			Name:  admin.Name,
		},
	}, nil
}

func (s *service) GetProfile(adminID uuid.UUID) (*AdminData, error) {
	admin, err := s.repo.FindByID(adminID)
	if err != nil {
		return nil, errors.New("admin profile not found")
	}

	return &AdminData{
		ID:    admin.ID,
		Email: admin.Email,
		Name:  admin.Name,
	}, nil
}
