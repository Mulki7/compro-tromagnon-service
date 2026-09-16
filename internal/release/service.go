package release

import (
	"errors"
	"strings"
	"time"

	"compro-tromagnon-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetPublishedReleases() ([]Release, error)
	GetAllReleases() ([]Release, error)
	GetBySlug(slug string) (*Release, error)
	GetByID(id uuid.UUID) (*Release, error)
	Create(req CreateReleaseRequest) (*Release, error)
	Update(id uuid.UUID, req UpdateReleaseRequest) (*Release, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPublishedReleases() ([]Release, error) {
	return s.repo.FindAll("published")
}

func (s *service) GetAllReleases() ([]Release, error) {
	return s.repo.FindAll("")
}

func (s *service) GetBySlug(slug string) (*Release, error) {
	rel, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}
	return rel, nil
}

func (s *service) GetByID(id uuid.UUID) (*Release, error) {
	rel, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}
	return rel, nil
}

func (s *service) Create(req CreateReleaseRequest) (*Release, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("release title is required")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = utils.MakeSlug(title)
	} else {
		slug = utils.MakeSlug(slug)
	}

	existing, _ := s.repo.FindBySlug(slug)
	if existing != nil {
		return nil, errors.New("a release with this slug already exists")
	}

	status := req.Status
	if status == "" {
		status = "published"
	}

	year := req.ReleaseYear
	if year == 0 {
		year = time.Now().Year()
	}

	rel := &Release{
		ArtistID:       req.ArtistID,
		Title:          title,
		Slug:           slug,
		ReleaseYear:    year,
		Format:         req.Format,
		ArtworkURL:     req.ArtworkURL,
		Description:    req.Description,
		StreamingLinks: req.StreamingLinks,
		Status:         status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(rel); err != nil {
		return nil, err
	}

	return s.repo.FindByID(rel.ID)
}

func (s *service) Update(id uuid.UUID, req UpdateReleaseRequest) (*Release, error) {
	rel, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title != "" {
			rel.Title = title
		}
	}

	if req.Slug != nil && *req.Slug != "" {
		newSlug := utils.MakeSlug(*req.Slug)
		if newSlug != rel.Slug {
			existing, _ := s.repo.FindBySlug(newSlug)
			if existing != nil && existing.ID != rel.ID {
				return nil, errors.New("a release with this slug already exists")
			}
			rel.Slug = newSlug
		}
	}

	if req.ArtistID != nil {
		rel.ArtistID = req.ArtistID
	}
	if req.ReleaseYear != nil && *req.ReleaseYear != 0 {
		rel.ReleaseYear = *req.ReleaseYear
	}
	if req.Format != nil {
		rel.Format = req.Format
	}
	if req.ArtworkURL != nil {
		rel.ArtworkURL = *req.ArtworkURL
	}
	if req.Description != nil {
		rel.Description = *req.Description
	}
	if req.StreamingLinks != nil {
		rel.StreamingLinks = req.StreamingLinks
	}
	if req.Status != nil && *req.Status != "" {
		rel.Status = *req.Status
	}
	rel.UpdatedAt = time.Now()

	if err := s.repo.Update(rel); err != nil {
		return nil, err
	}

	return s.repo.FindByID(rel.ID)
}

func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("release not found")
		}
		return err
	}
	return s.repo.Delete(id)
}
