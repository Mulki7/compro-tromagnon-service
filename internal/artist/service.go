package artist

import (
	"errors"
	"strings"
	"time"

	"compro-tromagnon-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetActiveArtists() ([]Artist, error)
	GetAllArtists() ([]Artist, error)
	GetBySlug(slug string) (*Artist, error)
	GetByID(id uuid.UUID) (*Artist, error)
	Create(req CreateArtistRequest) (*Artist, error)
	Update(id uuid.UUID, req UpdateArtistRequest) (*Artist, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetActiveArtists() ([]Artist, error) {
	return s.repo.FindAll("active")
}

func (s *service) GetAllArtists() ([]Artist, error) {
	return s.repo.FindAll("")
}

func (s *service) GetBySlug(slug string) (*Artist, error) {
	artist, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return artist, nil
}

func (s *service) GetByID(id uuid.UUID) (*Artist, error) {
	artist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return artist, nil
}

func (s *service) Create(req CreateArtistRequest) (*Artist, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("artist name is required")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = utils.MakeSlug(name)
	} else {
		slug = utils.MakeSlug(slug)
	}

	// Check if slug exists
	existing, _ := s.repo.FindBySlug(slug)
	if existing != nil {
		return nil, errors.New("an artist with this slug already exists")
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	artist := &Artist{
		Name:        name,
		Slug:        slug,
		Bio:         req.Bio,
		PhotoURL:    req.PhotoURL,
		Genre:       req.Genre,
		SocialLinks: req.SocialLinks,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(artist); err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *service) Update(id uuid.UUID, req UpdateArtistRequest) (*Artist, error) {
	artist, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name != "" {
			artist.Name = name
		}
	}

	if req.Slug != nil && *req.Slug != "" {
		newSlug := utils.MakeSlug(*req.Slug)
		if newSlug != artist.Slug {
			existing, _ := s.repo.FindBySlug(newSlug)
			if existing != nil && existing.ID != artist.ID {
				return nil, errors.New("an artist with this slug already exists")
			}
			artist.Slug = newSlug
		}
	}

	if req.Bio != nil {
		artist.Bio = *req.Bio
	}
	if req.PhotoURL != nil {
		artist.PhotoURL = *req.PhotoURL
	}
	if req.Genre != nil {
		artist.Genre = req.Genre
	}
	if req.SocialLinks != nil {
		artist.SocialLinks = req.SocialLinks
	}
	if req.Status != nil && *req.Status != "" {
		artist.Status = *req.Status
	}
	artist.UpdatedAt = time.Now()

	if err := s.repo.Update(artist); err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("artist not found")
		}
		return err
	}
	return s.repo.Delete(id)
}
