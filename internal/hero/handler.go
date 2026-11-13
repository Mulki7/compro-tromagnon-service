package hero

import (
	"strconv"

	"compro-tromagnon-service/internal/upload"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const uploadDir = "hero"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetBanners — Public: list semua banner urut by order
func (h *Handler) GetBanners(c *fiber.Ctx) error {
	banners, err := h.service.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve hero banners")
	}
	return response.Success(c, fiber.StatusOK, banners, "Hero banners retrieved successfully")
}

// GetBannerByID — Admin: detail banner by UUID
func (h *Handler) GetBannerByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid banner ID")
	}
	banner, err := h.service.GetByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, fiber.StatusOK, banner, "Hero banner retrieved successfully")
}

// CreateBanner — Admin: multipart/form-data
// Fields:
//
//	title         string (required)
//	subtitle      string
//	order         int    (urutan slideshow, default 0)
//	banner_image  file   (optional image upload)
func (h *Handler) CreateBanner(c *fiber.Ctx) error {
	imageURL, err := upload.SaveFromForm(c, "banner_image", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if imageURL == "" {
		imageURL = c.FormValue("image_url")
	}

	order := 0
	if o := c.FormValue("order"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			order = parsed
		}
	}

	req := CreateHeroBannerRequest{
		Title:    c.FormValue("title"),
		Subtitle: c.FormValue("subtitle"),
		ImageURL: imageURL,
		Order:    order,
	}

	banner, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, banner, "Hero banner created successfully")
}

// UpdateBanner — Admin: multipart/form-data (PUT & PATCH)
func (h *Handler) UpdateBanner(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid banner ID")
	}

	imageURL, err := upload.SaveFromForm(c, "banner_image", uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	var req UpdateHeroBannerRequest

	if title := c.FormValue("title"); title != "" {
		req.Title = &title
	}
	if subtitle := c.FormValue("subtitle"); subtitle != "" {
		req.Subtitle = &subtitle
	}
	if imageURL != "" {
		req.ImageURL = &imageURL
	} else if iURL := c.FormValue("image_url"); iURL != "" {
		req.ImageURL = &iURL
	}
	if o := c.FormValue("order"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			req.Order = &parsed
		}
	}

	banner, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, banner, "Hero banner updated successfully")
}

// DeleteBanner — Admin: soft delete
func (h *Handler) DeleteBanner(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid banner ID")
	}
	if err := h.service.Delete(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, fiber.StatusOK, nil, "Hero banner deleted successfully")
}
