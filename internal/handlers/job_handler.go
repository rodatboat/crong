package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
)

func CreateJob(c fiber.Ctx) error {
	newJob := new(models.Job)
	if err := c.Bind().Body(newJob); err != nil {
		return err
	}

	// TODO: Call repository to create job in database

	return response.Success(c, newJob)
}

func ReadJobs(c fiber.Ctx) error {

	// TODO: Call repository to read jobs from database

	jobs := []models.Job{}

	return response.Success(c, jobs)
}

func UpdateJob(c fiber.Ctx) error {

	jobId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid job ID")
	}

	// TODO: Call repository to update job in database

	job := models.Job{
		ID: uint(jobId),
	}

	return response.Success(c, &job)
}

func DeleteJob(c fiber.Ctx) error {
	_, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid job ID")
	}

	// TODO: Call repository to delete job from database

	return response.Success(c, nil)
}
