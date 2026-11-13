package live

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetUpcomingLiveEvents() ([]LiveEvent, error)
	GetAllLiveEvents() ([]LiveEvent, error)
	GetByID(id uuid.UUID) (*LiveEvent, error)
	Create(req CreateLiveEventRequest) (*LiveEvent, error)
	Update(id uuid.UUID, req UpdateLiveEventRequest) (*LiveEvent, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUpcomingLiveEvents() ([]LiveEvent, error) {
	return s.repo.FindAll("upcoming")
}

func (s *service) GetAllLiveEvents() ([]LiveEvent, error) {
	return s.repo.FindAll("")
}

func (s *service) GetByID(id uuid.UUID) (*LiveEvent, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("live event not found")
		}
		return nil, err
	}
	return event, nil
}

func (s *service) Create(req CreateLiveEventRequest) (*LiveEvent, error) {
	city := strings.TrimSpace(req.City)
	venue := strings.TrimSpace(req.Venue)
	if venue == "" {
		return nil, errors.New("venue name is required")
	}

	if req.EventDate == "" {
		return nil, errors.New("event_date is required (format: YYYY-MM-DD)")
	}

	parsedDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		return nil, errors.New("invalid event_date format, must be YYYY-MM-DD")
	}

	status := req.Status
	if status == "" {
		status = "upcoming"
	}

	country := req.Country
	if country == "" {
		country = "Indonesia"
	}

	event := &LiveEvent{
		ArtistID:  req.ArtistID,
		EventDate: parsedDate,
		City:      city,
		Venue:     venue,
		Country:   country,
		TicketURL: req.TicketURL,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}

	return s.repo.FindByID(event.ID)
}

func (s *service) Update(id uuid.UUID, req UpdateLiveEventRequest) (*LiveEvent, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("live event not found")
		}
		return nil, err
	}

	if req.ArtistID != nil {
		event.ArtistID = req.ArtistID
	}
	if req.EventDate != nil && *req.EventDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.EventDate)
		if err != nil {
			return nil, errors.New("invalid event_date format, must be YYYY-MM-DD")
		}
		event.EventDate = parsedDate
	}
	if req.City != nil {
		event.City = *req.City
	}
	if req.Venue != nil {
		event.Venue = *req.Venue
	}
	if req.Country != nil {
		event.Country = *req.Country
	}
	if req.TicketURL != nil {
		event.TicketURL = *req.TicketURL
	}
	if req.Status != nil && *req.Status != "" {
		event.Status = *req.Status
	}
	event.UpdatedAt = time.Now()

	if err := s.repo.Update(event); err != nil {
		return nil, err
	}

	return s.repo.FindByID(event.ID)
}

func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("live event not found")
		}
		return err
	}
	return s.repo.Delete(id)
}
