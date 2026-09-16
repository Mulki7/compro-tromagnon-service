package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*Admin, error)
	FindByID(id uuid.UUID) (*Admin, error)
	Create(admin *Admin) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(email string) (*Admin, error) {
	var admin Admin
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *repository) FindByID(id uuid.UUID) (*Admin, error) {
	var admin Admin
	err := r.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *repository) Create(admin *Admin) error {
	return r.db.Create(admin).Error
}
