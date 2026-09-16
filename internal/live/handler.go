package live

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

// Public Handlers

func (h *Handler) GetPublicLiveEvents(c *fiber.Ctx) error {
	events, err := h.service.GetUpcomingLiveEvents()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve live events")
	}
	return response.Success(c, fiber.StatusOK, events, "Live events retrieved successfully")
}

// Admin Handlers

func (h *Handler) GetAdminLiveEvents(c *fiber.Ctx) error {
	events, err := h.service.GetAllLiveEvents()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to retrieve live events")
	}
	return response.Success(c, fiber.StatusOK, events, "All live events retrieved successfully")
}

func (h *Handler) GetAdminLiveEventByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid live event ID")
	}

	event, err := h.service.GetByID(id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, event, "Live event retrieved successfully")
}

func (h *Handler) CreateLiveEvent(c *fiber.Ctx) error {
	var req CreateLiveEventRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	event, err := h.service.Create(req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, event, "Live event created successfully")
}

func (h *Handler) UpdateLiveEvent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid live event ID")
	}

	var req UpdateLiveEventRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	event, err := h.service.Update(id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, event, "Live event updated successfully")
}

func (h *Handler) DeleteLiveEvent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid live event ID")
	}

	if err := h.service.Delete(id); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil, "Live event deleted successfully")
}
