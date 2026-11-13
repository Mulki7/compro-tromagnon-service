package hero

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]HeroBanner, error)
	FindByID(id uuid.UUID) (*HeroBanner, error)
	Create(banner *HeroBanner) error
	Update(banner *HeroBanner) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindAll returns all active banners ordered by `order` ASC
func (r *repository) FindAll() ([]HeroBanner, error) {
	var banners []HeroBanner
	result := r.db.Order("\"order\" ASC, created_at ASC").Find(&banners)
	return banners, result.Error
}

func (r *repository) FindByID(id uuid.UUID) (*HeroBanner, error) {
	var banner HeroBanner
	result := r.db.First(&banner, "id = ?", id)
	return &banner, result.Error
}

func (r *repository) Create(banner *HeroBanner) error {
	return r.db.Create(banner).Error
}

func (r *repository) Update(banner *HeroBanner) error {
	return r.db.Save(banner).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&HeroBanner{}, "id = ?", id).Error
}
