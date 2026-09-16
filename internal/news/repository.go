package news

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(status string) ([]News, error)
	FindBySlug(slug string) (*News, error)
	FindByID(id uuid.UUID) (*News, error)
	Create(news *News) error
	Update(news *News) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(status string) ([]News, error) {
	var newsList []News
	query := r.db.Preload("RelatedArtist").Preload("RelatedRelease").Order("published_at DESC, created_at DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&newsList).Error
	return newsList, err
}

func (r *repository) FindBySlug(slug string) (*News, error) {
	var n News
	err := r.db.Preload("RelatedArtist").Preload("RelatedRelease").Where("slug = ?", slug).First(&n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *repository) FindByID(id uuid.UUID) (*News, error) {
	var n News
	err := r.db.Preload("RelatedArtist").Preload("RelatedRelease").Where("id = ?", id).First(&n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *repository) Create(news *News) error {
	return r.db.Create(news).Error
}

func (r *repository) Update(news *News) error {
	return r.db.Save(news).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&News{}, "id = ?", id).Error
}
