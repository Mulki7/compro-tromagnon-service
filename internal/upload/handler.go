package upload

import (
	"compro-tromagnon-service/internal/config"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg       *config.Config
	uploadDir string
}

func NewHandler(cfg *config.Config, uploadDir string) *Handler {
	if uploadDir == "" || uploadDir == "./uploads" {
		uploadDir = "general"
	}
	return &Handler{
		cfg:       cfg,
		uploadDir: uploadDir,
	}
}

// UploadFile — standalone endpoint POST /api/admin/upload
// Berguna jika ingin upload gambar terlebih dahulu sebelum mengisi form.
func (h *Handler) UploadFile(c *fiber.Ctx) error {
	fileURL, err := SaveFromForm(c, "file", h.uploadDir)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	// Try "image" field as fallback
	if fileURL == "" {
		fileURL, err = SaveFromForm(c, "image", h.uploadDir)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	}

	if fileURL == "" {
		return response.Error(c, fiber.StatusBadRequest, "No file uploaded. Use field 'file' or 'image'")
	}

	return response.Success(c, fiber.StatusCreated, fiber.Map{
		"url": fileURL,
	}, "File uploaded successfully")
}
