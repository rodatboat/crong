package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
	"github.com/rodatboat/crong/internal/services"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) CreateJob(c fiber.Ctx) error {
	var req models.JobCreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// TODO: Get user from context (from auth middleware)
	// userID := c.Locals("user_id").(uint)

	// Call service layer
	job, err := h.jobService.CreateJob(1, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, job)
}

func (h *JobHandler) ReadJobs(c fiber.Ctx) error {
	// TODO: Get user from context
	// userID := c.Locals("user_id").(uint)

	// Call service layer
	jobs, err := h.jobService.GetJobsByUser(1)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, jobs)
}

func (h *JobHandler) UpdateJob(c fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid job ID")
	}

	// TODO: Get user from context
	// userID := c.Locals("user_id").(uint)

	var req models.JobUpdateRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Call service layer
	job, err := h.jobService.UpdateJob(uint(jobID), 1, &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, job)
}

func (h *JobHandler) DeleteJob(c fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid job ID")
	}

	// TODO: Get user from context
	// userID := c.Locals("user_id").(uint)

	err = h.jobService.DeleteJob(uint(jobID), 1)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, nil)
}
