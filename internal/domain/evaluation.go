package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EvaluationStatus string

const (
	StatusQueued     EvaluationStatus = "queued"
	StatusProcessing EvaluationStatus = "processing"
	StatusCompleted  EvaluationStatus = "completed"
	StatusFailed     EvaluationStatus = "failed"
)

// Evaluation represents the core domain model
type Evaluation struct {
	ID         uuid.UUID        `db:"id"`
	Status     EvaluationStatus `db:"status"`
	CVPath     string           `db:"cv_path"`
	ReportPath string           `db:"report_path"`
	Result     *json.RawMessage `db:"result"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"`
}

// Struct for the final result format
type EvaluationResult struct {
	CVMatchRate     float64 `json:"cv_match_rate"`
	CVFeedback      string  `json:"cv_feedback"`
	ProjectScore    float64 `json:"project_score"`
	ProjectFeedback string  `json:"project_feedback"`
	OverallSummary  string  `json:"overall_summary"`
}
