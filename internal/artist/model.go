package artist

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Artist struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Bio         string         `gorm:"type:text" json:"bio"`
	PhotoURL    string         `gorm:"type:varchar(500)" json:"photo_url"`
	Genre       datatypes.JSON `gorm:"type:jsonb" json:"genre"`
	SocialLinks datatypes.JSON `gorm:"type:jsonb" json:"social_links"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Artist) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type CreateArtistRequest struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Bio         string         `json:"bio"`
	PhotoURL    string         `json:"photo_url"`
	Genre       datatypes.JSON `json:"genre"`
	SocialLinks datatypes.JSON `json:"social_links"`
	Status      string         `json:"status"`
}

type UpdateArtistRequest struct {
	Name        *string        `json:"name"`
	Slug        *string        `json:"slug"`
	Bio         *string        `json:"bio"`
	PhotoURL    *string        `json:"photo_url"`
	Genre       datatypes.JSON `json:"genre"`
	SocialLinks datatypes.JSON `json:"social_links"`
	Status      *string        `json:"status"`
}
