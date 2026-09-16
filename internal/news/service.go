package news

import (
	"errors"
	"strings"
	"time"

	"compro-tromagnon-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetPublishedNews() ([]News, error)
	GetAllNews() ([]News, error)
	GetBySlug(slug string) (*News, error)
	GetByID(id uuid.UUID) (*News, error)
	Create(req CreateNewsRequest) (*News, error)
	Update(id uuid.UUID, req UpdateNewsRequest) (*News, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPublishedNews() ([]News, error) {
	return s.repo.FindAll("published")
}

func (s *service) GetAllNews() ([]News, error) {
	return s.repo.FindAll("")
}

func (s *service) GetBySlug(slug string) (*News, error) {
	n, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}
	return n, nil
}

func (s *service) GetByID(id uuid.UUID) (*News, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}
	return n, nil
}

func (s *service) Create(req CreateNewsRequest) (*News, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("news title is required")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("news content is required")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = utils.MakeSlug(title)
	} else {
		slug = utils.MakeSlug(slug)
	}

	existing, _ := s.repo.FindBySlug(slug)
	if existing != nil {
		return nil, errors.New("a news article with this slug already exists")
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	publishedAt := req.PublishedAt
	if status == "published" && publishedAt == nil {
		now := time.Now()
		publishedAt = &now
	}

	newsItem := &News{
		Title:            title,
		Slug:             slug,
		Content:          content,
		CoverImageURL:    req.CoverImageURL,
		RelatedArtistID:  req.RelatedArtistID,
		RelatedReleaseID: req.RelatedReleaseID,
		PublishedAt:      publishedAt,
		Status:           status,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.Create(newsItem); err != nil {
		return nil, err
	}

	return s.repo.FindByID(newsItem.ID)
}

func (s *service) Update(id uuid.UUID, req UpdateNewsRequest) (*News, error) {
	newsItem, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title != "" {
			newsItem.Title = title
		}
	}

	if req.Content != nil {
		content := strings.TrimSpace(*req.Content)
		if content != "" {
			newsItem.Content = content
		}
	}

	if req.Slug != nil && *req.Slug != "" {
		newSlug := utils.MakeSlug(*req.Slug)
		if newSlug != newsItem.Slug {
			existing, _ := s.repo.FindBySlug(newSlug)
			if existing != nil && existing.ID != newsItem.ID {
				return nil, errors.New("a news article with this slug already exists")
			}
			newsItem.Slug = newSlug
		}
	}

	if req.CoverImageURL != nil {
		newsItem.CoverImageURL = *req.CoverImageURL
	}
	if req.RelatedArtistID != nil {
		newsItem.RelatedArtistID = req.RelatedArtistID
	}
	if req.RelatedReleaseID != nil {
		newsItem.RelatedReleaseID = req.RelatedReleaseID
	}
	if req.PublishedAt != nil {
		newsItem.PublishedAt = req.PublishedAt
	}
	if req.Status != nil && *req.Status != "" {
		if *req.Status == "published" && newsItem.PublishedAt == nil {
			now := time.Now()
			newsItem.PublishedAt = &now
		}
		newsItem.Status = *req.Status
	}
	newsItem.UpdatedAt = time.Now()

	if err := s.repo.Update(newsItem); err != nil {
		return nil, err
	}

	return s.repo.FindByID(newsItem.ID)
}

func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("news not found")
		}
		return err
	}
	return s.repo.Delete(id)
}
