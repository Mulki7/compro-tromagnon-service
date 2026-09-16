package release

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(status string) ([]Release, error)
	FindBySlug(slug string) (*Release, error)
	FindByID(id uuid.UUID) (*Release, error)
	Create(release *Release) error
	Update(release *Release) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(status string) ([]Release, error) {
	var releases []Release
	query := r.db.Preload("Artist").Order("release_year DESC, created_at DESC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&releases).Error
	return releases, err
}

func (r *repository) FindBySlug(slug string) (*Release, error) {
	var rel Release
	err := r.db.Preload("Artist").Where("slug = ?", slug).First(&rel).Error
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Release, error) {
	var rel Release
	err := r.db.Preload("Artist").Where("id = ?", id).First(&rel).Error
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *repository) Create(release *Release) error {
	return r.db.Create(release).Error
}

func (r *repository) Update(release *Release) error {
	return r.db.Save(release).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Release{}, "id = ?", id).Error
}
