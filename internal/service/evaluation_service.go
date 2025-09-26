package service

import (
	"context"
	"encoding/json"
	"log"
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
	go s.processEvaluation(eval.ID)

	return eval, nil
}

func (s *evaluationService) GetEvaluationResult(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error) {
	return s.repo.FindByID(ctx, id)
}

// TODO: Implement the processEvaluation method that will run in the background
func (s *evaluationService) processEvaluation(id uuid.UUID) {
	log.Printf("Starting evaluation for job ID: %s", id)

	// Simulasikan (misalnya, pemanggilan AI)
	time.Sleep(15 * time.Second)

	// Ambil data evaluasi dari DB
	eval, err := s.repo.FindByID(context.Background(), id)
	if err != nil {
		log.Printf("Error finding evaluation %s for processing: %v", id, err)
		// TODO: Update status ke 'failed'
		return
	}

	// TODO: logic pipeline AI
	// 1. Baca file CV dan Report dari path (eval.CVPath, eval.ReportPath)
	// 2. Lakukan RAG (ambil konteks dari ChromaDB)
	// 3. Lakukan LLM Chaining (panggil API Gemini)
	// 4. Hasilkan JSON result

	// dummy
	dummyResultJSON := `{"cv_match_rate": 0.82, "cv_feedback": "Strong in backend.", "project_score": 7.5, "project_feedback": "Good.", "overall_summary": "Promising candidate."}`

	// Update status dan hasil di database
	eval.Status = domain.StatusCompleted
	raw := json.RawMessage(dummyResultJSON)
	eval.Result = &raw

	err = s.repo.Update(context.Background(), eval)
	if err != nil {
		log.Printf("Error updating evaluation %s to completed: %v", id, err)
		// TODO: Update status ke 'failed'
		return
	}

	log.Printf("Finished evaluation for job ID: %s", id)
}
