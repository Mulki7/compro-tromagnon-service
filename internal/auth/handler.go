package auth

import (
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	res, err := h.service.Login(req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, res, "Login successful")
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	adminIDStr, ok := c.Locals("admin_id").(string)
	if !ok || adminIDStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid admin ID")
	}

	adminData, err := h.service.GetProfile(adminID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, adminData, "Profile retrieved successfully")
}
