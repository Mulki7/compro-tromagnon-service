package live

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(status string) ([]LiveEvent, error)
	FindByID(id uuid.UUID) (*LiveEvent, error)
	Create(event *LiveEvent) error
	Update(event *LiveEvent) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(status string) ([]LiveEvent, error) {
	var events []LiveEvent
	query := r.db.Preload("Artist")

	if status == "upcoming" {
		today := time.Now().Truncate(24 * time.Hour)
		query = query.Where("status = ? AND event_date >= ?", "upcoming", today).Order("event_date ASC")
	} else if status != "" {
		query = query.Where("status = ?", status).Order("event_date DESC")
	} else {
		query = query.Order("event_date DESC")
	}

	err := query.Find(&events).Error
	return events, err
}

func (r *repository) FindByID(id uuid.UUID) (*LiveEvent, error) {
	var event LiveEvent
	err := r.db.Preload("Artist").Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *repository) Create(event *LiveEvent) error {
	return r.db.Create(event).Error
}

func (r *repository) Update(event *LiveEvent) error {
	return r.db.Save(event).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&LiveEvent{}, "id = ?", id).Error
}
