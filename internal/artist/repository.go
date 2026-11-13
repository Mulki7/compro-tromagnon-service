package artist

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(status string) ([]Artist, error)
	FindBySlug(slug string) (*Artist, error)
	FindByID(id uuid.UUID) (*Artist, error)
	Create(artist *Artist) error
	Update(artist *Artist) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(status string) ([]Artist, error) {
	var artists []Artist
	query := r.db.Order("name ASC")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&artists).Error
	return artists, err
}

func (r *repository) FindBySlug(slug string) (*Artist, error) {
	var artist Artist
	err := r.db.Where("slug = ?", slug).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Artist, error) {
	var artist Artist
	err := r.db.Where("id = ?", id).First(&artist).Error
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

func (r *repository) Create(artist *Artist) error {
	return r.db.Create(artist).Error
}

func (r *repository) Update(artist *Artist) error {
	return r.db.Save(artist).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Artist{}, "id = ?", id).Error
}
