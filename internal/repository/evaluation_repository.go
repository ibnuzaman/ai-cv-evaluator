package repository

import (
	"context"

	"aicvevaluator/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// EvaluationRepository defines the contract for database operations
type EvaluationRepository interface {
	Create(ctx context.Context, evaluation *domain.Evaluation) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error)
	Update(ctx context.Context, evaluation *domain.Evaluation) error
}

// postgresEvaluationRepo implements EvaluationRepository for PostgreSQL
type postgresEvaluationRepo struct {
	db *sqlx.DB
}

// NewEvaluationRepository creates a new instance of the repository
func NewEvaluationRepository(db *sqlx.DB) EvaluationRepository {
	return &postgresEvaluationRepo{db: db}
}

func (r *postgresEvaluationRepo) Create(ctx context.Context, eval *domain.Evaluation) error {
	query := `INSERT INTO evaluations (id, status, cv_path, report_path, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, eval.ID, eval.Status, eval.CVPath, eval.ReportPath, eval.CreatedAt, eval.UpdatedAt)
	return err
}

func (r *postgresEvaluationRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Evaluation, error) {
	var eval domain.Evaluation
	query := `SELECT id, status, cv_path, report_path, result, created_at, updated_at
			  FROM evaluations WHERE id = $1`
	err := r.db.GetContext(ctx, &eval, query, id)
	return &eval, err
}

func (r *postgresEvaluationRepo) Update(ctx context.Context, eval *domain.Evaluation) error {
	query := `UPDATE evaluations 
			  SET status = $2, result = $3, updated_at = NOW()
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, eval.ID, eval.Status, eval.Result)
	return err
}
