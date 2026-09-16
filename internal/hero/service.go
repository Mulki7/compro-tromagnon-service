package hero

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAll() ([]HeroBanner, error)
	GetByID(id uuid.UUID) (*HeroBanner, error)
	Create(req CreateHeroBannerRequest) (*HeroBanner, error)
	Update(id uuid.UUID, req UpdateHeroBannerRequest) (*HeroBanner, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]HeroBanner, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uuid.UUID) (*HeroBanner, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hero banner not found")
		}
		return nil, err
	}
	return banner, nil
}

func (s *service) Create(req CreateHeroBannerRequest) (*HeroBanner, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("banner title is required")
	}

	banner := &HeroBanner{
		Title:     title,
		Subtitle:  req.Subtitle,
		ImageURL:  req.ImageURL,
		Order:     req.Order,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(banner); err != nil {
		return nil, err
	}

	return banner, nil
}

func (s *service) Update(id uuid.UUID, req UpdateHeroBannerRequest) (*HeroBanner, error) {
	banner, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hero banner not found")
		}
		return nil, err
	}

	if req.Title != nil && *req.Title != "" {
		banner.Title = strings.TrimSpace(*req.Title)
	}
	if req.Subtitle != nil {
		banner.Subtitle = *req.Subtitle
	}
	if req.ImageURL != nil {
		banner.ImageURL = *req.ImageURL
	}
	if req.Order != nil {
		banner.Order = *req.Order
	}
	banner.UpdatedAt = time.Now()

	if err := s.repo.Update(banner); err != nil {
		return nil, err
	}

	return banner, nil
}

func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("hero banner not found")
		}
		return err
	}
	return s.repo.Delete(id)
}
