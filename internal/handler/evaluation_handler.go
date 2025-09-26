package handler

import (
	"aicvevaluator/internal/service"
	"log"

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
	// TODO: Handle file uploads from the request
	// For now, we use placeholders
	cvPath := "path/to/cv.pdf"
	reportPath := "path/to/report.pdf"

	eval, err := h.service.CreateEvaluation(c.Context(), cvPath, reportPath)
	if err != nil {
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
