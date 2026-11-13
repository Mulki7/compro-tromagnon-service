package release

import (
	"encoding/json"
	"strconv"

	"compro-tromagnon-service/internal/upload"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const uploadDir = "release"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Public Handlers
// =============================================================================

func (h *Handler) GetPublicReleases(c *fiber.Ctx) error {
	releases, err := h.service.GetPublishedReleases()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve releases")
	}
	return response.Success(c, fiber.StatusOK, releases, "Releases retrieved successfully")
}

func (h *Handler) GetPublicReleaseBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "Slug parameter is required")
	}
	rel, err := h.service.GetBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, rel, "Release retrieved successfully")
}

// =============================================================================
// Admin Handlers
// =============================================================================

func (h *Handler) GetAdminReleases(c *fiber.Ctx) error {
	releases, err := h.service.GetAllReleases()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve releases")
	}
	return response.Success(c, fiber.StatusOK, releases, "All releases retrieved successfully")
}

func (h *Handler) GetAdminReleaseByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid release ID")
	}
	rel, err := h.service.GetByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, rel, "Release retrieved successfully")
}

// CreateRelease — multipart/form-data
// Fields:
//
//	artist_id       uuid string
//	title           string (required)
//	slug            string (optional)
//	release_year    int
//	format          string (LP/EP/CD/Tape/Digital/Single)
//	description     string
//	streaming_links JSON string (e.g. {"spotify":"...","apple_music":"..."})
//	status          string (published/draft)
//	artwork         file   (optional image upload)
func (h *Handler) CreateRelease(c *fiber.Ctx) error {
	artworkURL, err := upload.SaveFromForm(c, "artwork", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if artworkURL == "" {
		artworkURL = c.FormValue("artwork_url")
	}

	var streamingLinks datatypes.JSON
	if sl := c.FormValue("streaming_links"); sl != "" {
		if err := json.Unmarshal([]byte(sl), &streamingLinks); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid streaming_links JSON format")
		}
	}

	var artistID *uuid.UUID
	if aid := c.FormValue("artist_id"); aid != "" {
		parsed, err := uuid.Parse(aid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid artist_id format")
		}
		artistID = &parsed
	}

	releaseYear := 0
	if ry := c.FormValue("release_year"); ry != "" {
		if y, err := strconv.Atoi(ry); err == nil {
			releaseYear = y
		}
	}

	var formatJSON datatypes.JSON
	if f := c.FormValue("format"); f != "" {
		if err := json.Unmarshal([]byte(f), &formatJSON); err != nil {
			// Wrap plain string menjadi array
			formatJSON = datatypes.JSON([]byte(`["` + f + `"]`))
		}
	}

	req := CreateReleaseRequest{
		ArtistID:       artistID,
		Title:          c.FormValue("title"),
		Slug:           c.FormValue("slug"),
		ReleaseYear:    releaseYear,
		Format:         formatJSON,
		ArtworkURL:     artworkURL,
		Description:    c.FormValue("description"),
		StreamingLinks: streamingLinks,
		Status:         c.FormValue("status"),
	}

	rel, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, rel, "Release created successfully")
}

// UpdateRelease — multipart/form-data (PUT & PATCH)
func (h *Handler) UpdateRelease(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid release ID")
	}

	artworkURL, err := upload.SaveFromForm(c, "artwork", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var req UpdateReleaseRequest

	if aid := c.FormValue("artist_id"); aid != "" {
		parsed, err := uuid.Parse(aid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid artist_id format")
		}
		req.ArtistID = &parsed
	}
	if title := c.FormValue("title"); title != "" {
		req.Title = &title
	}
	if slug := c.FormValue("slug"); slug != "" {
		req.Slug = &slug
	}
	if ry := c.FormValue("release_year"); ry != "" {
		if y, err := strconv.Atoi(ry); err == nil {
			req.ReleaseYear = &y
		}
	}
	if format := c.FormValue("format"); format != "" {
		var formatJSON datatypes.JSON
		if err := json.Unmarshal([]byte(format), &formatJSON); err != nil {
			// Wrap plain string menjadi array
			formatJSON = datatypes.JSON([]byte(`["` + format + `"]`))
		}
		req.Format = formatJSON
	}
	if desc := c.FormValue("description"); desc != "" {
		req.Description = &desc
	}
	if status := c.FormValue("status"); status != "" {
		req.Status = &status
	}
	if artworkURL != "" {
		req.ArtworkURL = &artworkURL
	} else if aURL := c.FormValue("artwork_url"); aURL != "" {
		req.ArtworkURL = &aURL
	}
	if sl := c.FormValue("streaming_links"); sl != "" {
		var streamingLinks datatypes.JSON
		if err := json.Unmarshal([]byte(sl), &streamingLinks); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid streaming_links JSON format")
		}
		req.StreamingLinks = streamingLinks
	}

	rel, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, rel, "Release updated successfully")
}

func (h *Handler) DeleteRelease(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid release ID")
	}
	if err := h.service.Delete(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, nil, "Release deleted successfully")
}
