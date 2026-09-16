package news

import (
	"compro-tromagnon-service/internal/upload"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const uploadDir = "news"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Public Handlers
// =============================================================================

func (h *Handler) GetPublicNews(c *fiber.Ctx) error {
	articles, err := h.service.GetPublishedNews()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve news")
	}
	return response.Success(c, fiber.StatusOK, articles, "News retrieved successfully")
}

func (h *Handler) GetPublicNewsBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "Slug parameter is required")
	}
	article, err := h.service.GetBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, article, "News retrieved successfully")
}

// =============================================================================
// Admin Handlers
// =============================================================================

func (h *Handler) GetAdminNews(c *fiber.Ctx) error {
	articles, err := h.service.GetAllNews()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve news")
	}
	return response.Success(c, fiber.StatusOK, articles, "All news retrieved successfully")
}

func (h *Handler) GetAdminNewsByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid news ID")
	}
	article, err := h.service.GetByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, article, "News retrieved successfully")
}

// CreateNews — multipart/form-data
// Fields:
//
//	title              string (required)
//	slug               string (optional)
//	content            string
//	status             string (published/draft)
//	related_artist_id  uuid string (optional)
//	related_release_id uuid string (optional)
//	cover              file   (optional image upload)
func (h *Handler) CreateNews(c *fiber.Ctx) error {
	coverURL, err := upload.SaveFromForm(c, "cover", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if coverURL == "" {
		coverURL = c.FormValue("cover_image_url")
	}

	req := CreateNewsRequest{
		Title:          c.FormValue("title"),
		Slug:           c.FormValue("slug"),
		Content:        c.FormValue("content"),
		CoverImageURL:  coverURL,
		Status:         c.FormValue("status"),
	}

	if aid := c.FormValue("related_artist_id"); aid != "" {
		parsed, err := uuid.Parse(aid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid related_artist_id format")
		}
		req.RelatedArtistID = &parsed
	}
	if rid := c.FormValue("related_release_id"); rid != "" {
		parsed, err := uuid.Parse(rid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid related_release_id format")
		}
		req.RelatedReleaseID = &parsed
	}

	article, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, article, "News created successfully")
}

// UpdateNews — multipart/form-data (PUT & PATCH)
func (h *Handler) UpdateNews(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid news ID")
	}

	coverURL, err := upload.SaveFromForm(c, "cover", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var req UpdateNewsRequest

	if title := c.FormValue("title"); title != "" {
		req.Title = &title
	}
	if slug := c.FormValue("slug"); slug != "" {
		req.Slug = &slug
	}
	if content := c.FormValue("content"); content != "" {
		req.Content = &content
	}
	if status := c.FormValue("status"); status != "" {
		req.Status = &status
	}
	if coverURL != "" {
		req.CoverImageURL = &coverURL
	} else if cURL := c.FormValue("cover_image_url"); cURL != "" {
		req.CoverImageURL = &cURL
	}
	if aid := c.FormValue("related_artist_id"); aid != "" {
		parsed, err := uuid.Parse(aid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid related_artist_id format")
		}
		req.RelatedArtistID = &parsed
	}
	if rid := c.FormValue("related_release_id"); rid != "" {
		parsed, err := uuid.Parse(rid)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid related_release_id format")
		}
		req.RelatedReleaseID = &parsed
	}

	article, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, article, "News updated successfully")
}

func (h *Handler) DeleteNews(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid news ID")
	}
	if err := h.service.Delete(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, nil, "News deleted successfully")
}
