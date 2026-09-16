package news

import (
	"time"

	"compro-tromagnon-service/internal/artist"
	"compro-tromagnon-service/internal/release"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type News struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	Title            string           `gorm:"type:varchar(255);not null" json:"title"`
	Slug             string           `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content          string           `gorm:"type:text;not null" json:"content"`
	CoverImageURL    string           `gorm:"type:varchar(500)" json:"cover_image_url"`
	RelatedArtistID  *uuid.UUID       `gorm:"type:uuid" json:"related_artist_id"`
	RelatedArtist    *artist.Artist   `gorm:"foreignKey:RelatedArtistID;constraint:OnDelete:SET NULL" json:"related_artist,omitempty"`
	RelatedReleaseID *uuid.UUID       `gorm:"type:uuid" json:"related_release_id"`
	RelatedRelease   *release.Release `gorm:"foreignKey:RelatedReleaseID;constraint:OnDelete:SET NULL" json:"related_release,omitempty"`
	PublishedAt      *time.Time       `gorm:"type:timestamptz" json:"published_at"`
	Status           string           `gorm:"type:varchar(20);default:'draft'" json:"status"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (n *News) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type CreateNewsRequest struct {
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	Content          string     `json:"content"`
	CoverImageURL    string     `json:"cover_image_url"`
	RelatedArtistID  *uuid.UUID `json:"related_artist_id"`
	RelatedReleaseID *uuid.UUID `json:"related_release_id"`
	PublishedAt      *time.Time `json:"published_at"`
	Status           string     `json:"status"`
}

type UpdateNewsRequest struct {
	Title            *string    `json:"title"`
	Slug             *string    `json:"slug"`
	Content          *string    `json:"content"`
	CoverImageURL    *string    `json:"cover_image_url"`
	RelatedArtistID  *uuid.UUID `json:"related_artist_id"`
	RelatedReleaseID *uuid.UUID `json:"related_release_id"`
	PublishedAt      *time.Time `json:"published_at"`
	Status           *string    `json:"status"`
}
