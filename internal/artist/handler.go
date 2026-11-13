package artist

import (
	"encoding/json"

	"compro-tromagnon-service/internal/upload"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const uploadDir = "artist"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Public Handlers
// =============================================================================

func (h *Handler) GetPublicArtists(c *fiber.Ctx) error {
	artists, err := h.service.GetActiveArtists()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve artists")
	}
	return response.Success(c, fiber.StatusOK, artists, "Artists retrieved successfully")
}

func (h *Handler) GetPublicArtistBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "Slug parameter is required")
	}
	artist, err := h.service.GetBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, artist, "Artist retrieved successfully")
}

// =============================================================================
// Admin Handlers
// =============================================================================

func (h *Handler) GetAdminArtists(c *fiber.Ctx) error {
	artists, err := h.service.GetAllArtists()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve artists")
	}
	return response.Success(c, fiber.StatusOK, artists, "All artists retrieved successfully")
}

func (h *Handler) GetAdminArtistByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid artist ID")
	}
	artist, err := h.service.GetByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, artist, "Artist retrieved successfully")
}

// CreateArtist — menerima multipart/form-data
// Fields:
//
//	name         string (required)
//	slug         string (optional)
//	bio          string
//	genre        string
//	social_links JSON string (e.g. {"instagram":"...","spotify":"..."})
//	status       string (active/inactive/draft)
//	photo        file   (optional image upload)
func (h *Handler) CreateArtist(c *fiber.Ctx) error {
	// Coba simpan foto jika ada
	photoURL, err := upload.SaveFromForm(c, "photo", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	// Jika tidak ada file upload, cek apakah ada field photo_url (URL manual)
	if photoURL == "" {
		photoURL = c.FormValue("photo_url")
	}

	// Parse genre dari JSON array string
	var genre datatypes.JSON
	if g := c.FormValue("genre"); g != "" {
		if err := json.Unmarshal([]byte(g), &genre); err != nil {
			// Bisa juga dikirim sebagai plain string, wrap jadi array
			genre = datatypes.JSON([]byte(`["` + g + `"]`))
		}
	}

	// Parse social_links dari JSON string
	var socialLinks datatypes.JSON
	if sl := c.FormValue("social_links"); sl != "" {
		if err := json.Unmarshal([]byte(sl), &socialLinks); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid social_links JSON format")
		}
	}

	req := CreateArtistRequest{
		Name:        c.FormValue("name"),
		Slug:        c.FormValue("slug"),
		Bio:         c.FormValue("bio"),
		PhotoURL:    photoURL,
		Genre:       genre,
		SocialLinks: socialLinks,
		Status:      c.FormValue("status"),
	}

	artist, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, artist, "Artist created successfully")
}

// UpdateArtist — menerima multipart/form-data (PUT & PATCH)
// Field yang tidak dikirim tidak akan diubah (pointer-based partial update)
func (h *Handler) UpdateArtist(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid artist ID")
	}

	// Upload foto baru jika ada
	photoURL, err := upload.SaveFromForm(c, "photo", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var req UpdateArtistRequest

	// Assign pointer hanya jika field dikirim
	if name := c.FormValue("name"); name != "" {
		req.Name = &name
	}
	if slug := c.FormValue("slug"); slug != "" {
		req.Slug = &slug
	}
	if bio := c.FormValue("bio"); bio != "" {
		req.Bio = &bio
	}
	if genre := c.FormValue("genre"); genre != "" {
		var genreJSON datatypes.JSON
		if err := json.Unmarshal([]byte(genre), &genreJSON); err != nil {
			// Wrap plain string menjadi array
			genreJSON = datatypes.JSON([]byte(`["` + genre + `"]`))
		}
		req.Genre = genreJSON
	}
	if status := c.FormValue("status"); status != "" {
		req.Status = &status
	}

	// photo: dari file upload atau field photo_url
	if photoURL != "" {
		req.PhotoURL = &photoURL
	} else if pURL := c.FormValue("photo_url"); pURL != "" {
		req.PhotoURL = &pURL
	}

	// social_links
	if sl := c.FormValue("social_links"); sl != "" {
		var socialLinks datatypes.JSON
		if err := json.Unmarshal([]byte(sl), &socialLinks); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid social_links JSON format")
		}
		req.SocialLinks = socialLinks
	}

	artist, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, artist, "Artist updated successfully")
}

func (h *Handler) DeleteArtist(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid artist ID")
	}
	if err := h.service.Delete(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, nil, "Artist deleted successfully")
}
