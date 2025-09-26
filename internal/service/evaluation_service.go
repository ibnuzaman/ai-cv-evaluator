package service

import (
	"context"
	"time"

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
	repo repository.EvaluationRepository
}

// NewEvaluationService creates a new instance of the service
func NewEvaluationService(repo repository.EvaluationRepository) EvaluationService {
	return &evaluationService{repo: repo}
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
	// go s.processEvaluation(eval.ID)

	return eval, nil
}

func (s *evaluationService) GetEvaluationResult(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error) {
	return s.repo.FindByID(ctx, id)
}

// TODO: Implement the processEvaluation method that will run in the background
// func (s *evaluationService) processEvaluation(id uuid.UUID) { ... }
