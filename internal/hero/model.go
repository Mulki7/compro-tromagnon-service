package hero

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HeroBanner — satu record = satu banner/slide di homepage
type HeroBanner struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Subtitle  string         `gorm:"type:text" json:"subtitle"`
	ImageURL  string         `gorm:"type:varchar(500)" json:"image_url"`
	Order     int            `gorm:"type:integer;default:0" json:"order"` // urutan tampil di slideshow
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (h *HeroBanner) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

type CreateHeroBannerRequest struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"image_url"`
	Order    int    `json:"order"`
}

type UpdateHeroBannerRequest struct {
	Title    *string `json:"title"`
	Subtitle *string `json:"subtitle"`
	ImageURL *string `json:"image_url"`
	Order    *int    `json:"order"`
}
