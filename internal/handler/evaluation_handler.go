package handler

import (
	"aicvevaluator/internal/service"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EvaluationHandler struct {
	service service.EvaluationService
}

func NewEvaluationHandler(s service.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{service: s}
}

func (h *EvaluationHandler) Evaluate(c *fiber.Ctx) error {
	cvFile, err := c.FormFile("cv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CV file is required"})
	}

	reportFile, err := c.FormFile("project_report")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project report file is required"})
	}

	cvFilename := fmt.Sprintf("%s-%s", uuid.New().String(), filepath.Base(cvFile.Filename))
	reportFilename := fmt.Sprintf("%s-%s", uuid.New().String(), filepath.Base(reportFile.Filename))

	cvPath := filepath.Join("uploads", cvFilename)
	reportPath := filepath.Join("uploads", reportFilename)

	if err := c.SaveFile(cvFile, cvPath); err != nil {
		log.Printf("Error saving CV file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save CV file"})
	}
	if err := c.SaveFile(reportFile, reportPath); err != nil {
		log.Printf("Error saving report file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save report file"})
	}

	eval, err := h.service.CreateEvaluation(c.Context(), cvPath, reportPath)
	if err != nil {
		log.Printf("Error creating evaluation task: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create evaluation task",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"id":     eval.ID.String(),
		"status": eval.Status,
	})
}

func (h *EvaluationHandler) GetResult(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	result, err := h.service.GetEvaluationResult(c.Context(), id)
	if err != nil {
		log.Printf("Error getting result for ID %s: %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "result not found"})
	}

	// Prepare response based on status
	response := fiber.Map{
		"id":     result.ID.String(),
		"status": result.Status,
	}

	if result.Status == "completed" {
		response["result"] = result.Result
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
