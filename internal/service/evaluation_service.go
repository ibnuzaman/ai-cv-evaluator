package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"aicvevaluator/internal/ai"
	"aicvevaluator/internal/domain"
	"aicvevaluator/internal/repository"

	"github.com/google/uuid"
)

// EvaluationService defines the business logic operations
type EvaluationService interface {
	CreateEvaluation(ctx context.Context, cvPath, reportPath string) (*domain.Evaluation, error)
	GetEvaluationResult(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error)
}

type evaluationService struct {
	repo       repository.EvaluationRepository
	aiPipeline *ai.Pipeline
}

// NewEvaluationService creates a new instance of the service
func NewEvaluationService(repo repository.EvaluationRepository, aiPipeline *ai.Pipeline) EvaluationService {
	return &evaluationService{
		repo:       repo,
		aiPipeline: aiPipeline,
	}
}

func (s *evaluationService) CreateEvaluation(ctx context.Context, cvPath, reportPath string) (*domain.Evaluation, error) {
	eval := &domain.Evaluation{
		ID:         uuid.New(),
		Status:     domain.StatusQueued,
		CVPath:     cvPath,     // Placeholder, will be filled from upload
		ReportPath: reportPath, // Placeholder
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := s.repo.Create(ctx, eval)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger the background processing (AI Pipeline) in a goroutine here
	go s.processEvaluation(eval.ID)

	return eval, nil
}

func (s *evaluationService) GetEvaluationResult(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error) {
	return s.repo.FindByID(ctx, id)
}

// processEvaluation runs the AI evaluation pipeline in the background
func (s *evaluationService) processEvaluation(id uuid.UUID) {
	log.Printf("Starting AI evaluation for job ID: %s", id)

	ctx := context.Background()

	// Update status to processing
	eval, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("Error finding evaluation %s for processing: %v", id, err)
		return
	}

	eval.Status = domain.StatusProcessing
	err = s.repo.Update(ctx, eval)
	if err != nil {
		log.Printf("Error updating evaluation %s to processing: %v", id, err)
		return
	}

	// Run the AI pipeline
	result, err := s.aiPipeline.ProcessEvaluation(ctx, eval.CVPath, eval.ReportPath)
	if err != nil {
		log.Printf("AI pipeline failed for evaluation %s: %v", id, err)

		// Update status to failed
		eval.Status = domain.StatusFailed
		s.repo.Update(ctx, eval)
		return
	}

	// Convert result to JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshaling result for evaluation %s: %v", id, err)

		// Update status to failed
		eval.Status = domain.StatusFailed
		s.repo.Update(ctx, eval)
		return
	}

	// Update evaluation with results
	eval.Status = domain.StatusCompleted
	raw := json.RawMessage(resultJSON)
	eval.Result = &raw

	err = s.repo.Update(ctx, eval)
	if err != nil {
		log.Printf("Error updating evaluation %s to completed: %v", id, err)
		return
	}

	log.Printf("Successfully completed AI evaluation for job ID: %s", id)
}
