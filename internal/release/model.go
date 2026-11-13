package release

import (
	"time"

	"compro-tromagnon-service/internal/artist"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Release struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ArtistID       *uuid.UUID     `gorm:"type:uuid" json:"artist_id"`
	Artist         *artist.Artist `gorm:"foreignKey:ArtistID;constraint:OnDelete:SET NULL" json:"artist,omitempty"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	ReleaseYear    int            `gorm:"type:integer" json:"release_year"`
	Format         datatypes.JSON `gorm:"type:jsonb" json:"format"`
	ArtworkURL     string         `gorm:"type:varchar(500)" json:"artwork_url"`
	Description    string         `gorm:"type:text" json:"description"`
	StreamingLinks datatypes.JSON `gorm:"type:jsonb" json:"streaming_links"`
	Status         string         `gorm:"type:varchar(20);default:'published'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Release) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type CreateReleaseRequest struct {
	ArtistID       *uuid.UUID     `json:"artist_id"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug"`
	ReleaseYear    int            `json:"release_year"`
	Format         datatypes.JSON `json:"format"`
	ArtworkURL     string         `json:"artwork_url"`
	Description    string         `json:"description"`
	StreamingLinks datatypes.JSON `json:"streaming_links"`
	Status         string         `json:"status"`
}

type UpdateReleaseRequest struct {
	ArtistID       *uuid.UUID     `json:"artist_id"`
	Title          *string        `json:"title"`
	Slug           *string        `json:"slug"`
	ReleaseYear    *int           `json:"release_year"`
	Format         datatypes.JSON `json:"format"`
	ArtworkURL     *string        `json:"artwork_url"`
	Description    *string        `json:"description"`
	StreamingLinks datatypes.JSON `json:"streaming_links"`
	Status         *string        `json:"status"`
}
