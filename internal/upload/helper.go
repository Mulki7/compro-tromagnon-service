package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	BaseUploadDir = "./uploads"
	MaxFileSize   = 5 * 1024 * 1024 // 5 MB
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
	".svg":  true,
}

// SaveFromForm tries to extract and save a file from multipart form into a specific subfolder of uploads.
// subFolder examples: "artist", "release", "news", "hero", "general".
// Returns the public URL string (e.g. "/uploads/artist/xxx.jpg"), empty string if no file uploaded, or error.
func SaveFromForm(c *fiber.Ctx, fieldName string, subFolder string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		// No file in this field — not an error, just no upload
		return "", nil
	}

	if file.Size == 0 {
		return "", nil
	}

	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedExtensions[ext] {
		return "", fmt.Errorf("invalid file format. Allowed: jpg, jpeg, png, webp, gif, svg")
	}

	// Normalize subFolder: e.g. "artist", or if someone passes "./uploads/artist", handle cleanly
	cleanSub := strings.TrimPrefix(subFolder, "./uploads")
	cleanSub = strings.TrimPrefix(cleanSub, "uploads")
	cleanSub = strings.Trim(cleanSub, "/\\")

	targetDir := BaseUploadDir
	if cleanSub != "" {
		targetDir = filepath.Join(BaseUploadDir, cleanSub)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to prepare upload directory: %w", err)
	}

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	dst := filepath.Join(targetDir, filename)

	if err := c.SaveFile(file, dst); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	if cleanSub != "" {
		return fmt.Sprintf("/uploads/%s/%s", cleanSub, filename), nil
	}
	return fmt.Sprintf("/uploads/%s", filename), nil
}
