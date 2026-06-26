package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/middleware"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/resp"
	"github.com/rodatboat/crong/internal/services"
	"github.com/rodatboat/crong/internal/utils"
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
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	auth := c.Locals(middleware.AuthContextKey).(middleware.AuthContext)

	// Call service layer
	job, err := h.jobService.CreateJob(auth.UserID, &req)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(job))
}

func (h *JobHandler) ReadJobs(c fiber.Ctx) error {
	auth := c.Locals(middleware.AuthContextKey).(middleware.AuthContext)

	// Call service layer
	jobs, err := h.jobService.GetJobsByUser(auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(jobs))
}

func (h *JobHandler) GetJobsDetailsByID(c fiber.Ctx) error {
	auth := c.Locals(middleware.AuthContextKey).(middleware.AuthContext)

	folderIDStr := c.Params("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil || folderID == 0 {
		return resp.Send(c, resp.BadRequest())
	}

	// Call service layer
	jobs, err := h.jobService.GetJobsDetailsByID(uint(folderID), auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(jobs))
}

func (h *JobHandler) UpdateJob(c fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 32)
	if err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	auth := c.Locals(middleware.AuthContextKey).(middleware.AuthContext)

	var req models.JobUpdateRequest
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	// Update job details
	job, err := h.jobService.UpdateJob(uint(jobID), auth.UserID, &req)
	if err != nil {
		return resp.Send(c, resp.InternalServerError())
	}

	return resp.Send(c, resp.Success(job))
}

func (h *JobHandler) DeleteJob(c fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 32)
	if err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	auth := c.Locals(middleware.AuthContextKey).(middleware.AuthContext)

	err = h.jobService.DeleteJob(uint(jobID), auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(nil))
}
