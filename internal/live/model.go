package live

import (
	"time"

	"compro-tromagnon-service/internal/artist"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LiveEvent struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ArtistID  *uuid.UUID     `gorm:"type:uuid" json:"artist_id"`
	Artist    *artist.Artist `gorm:"foreignKey:ArtistID;constraint:OnDelete:SET NULL" json:"artist,omitempty"`
	EventDate time.Time      `gorm:"type:date;not null" json:"event_date"`
	City      string         `gorm:"type:varchar(100)" json:"city"`
	Venue     string         `gorm:"type:varchar(255)" json:"venue"`
	Country   string         `gorm:"type:varchar(100)" json:"country"`
	TicketURL string         `gorm:"type:varchar(500)" json:"ticket_url"`
	Status    string         `gorm:"type:varchar(20);default:'upcoming'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l *LiveEvent) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type CreateLiveEventRequest struct {
	ArtistID  *uuid.UUID `json:"artist_id"`
	EventDate string     `json:"event_date"` // YYYY-MM-DD format
	City      string     `json:"city"`
	Venue     string     `json:"venue"`
	Country   string     `json:"country"`
	TicketURL string     `json:"ticket_url"`
	Status    string     `json:"status"`
}

type UpdateLiveEventRequest struct {
	ArtistID  *uuid.UUID `json:"artist_id"`
	EventDate *string    `json:"event_date"`
	City      *string    `json:"city"`
	Venue     *string    `json:"venue"`
	Country   *string    `json:"country"`
	TicketURL *string    `json:"ticket_url"`
	Status    *string    `json:"status"`
}
