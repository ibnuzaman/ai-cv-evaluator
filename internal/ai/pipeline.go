package ai

import (
	"aicvevaluator/internal/chromadb"
	"aicvevaluator/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Pipeline orchestrates the AI evaluation process
type Pipeline struct {
	fileReader   *util.FileReader
	chromaClient *chromadb.Client
	geminiClient *GeminiClient
}

// NewPipeline creates a new AI pipeline
func NewPipeline(fileReader *util.FileReader, chromaClient *chromadb.Client, geminiClient *GeminiClient) *Pipeline {
	return &Pipeline{
		fileReader:   fileReader,
		chromaClient: chromaClient,
		geminiClient: geminiClient,
	}
}

// EvaluationResult represents the final evaluation result
type EvaluationResult struct {
	CVMatchRate     float64 `json:"cv_match_rate"`
	CVFeedback      string  `json:"cv_feedback"`
	ProjectScore    float64 `json:"project_score"`
	ProjectFeedback string  `json:"project_feedback"`
	OverallSummary  string  `json:"overall_summary"`
}

// ProcessEvaluation runs the complete AI evaluation pipeline
func (p *Pipeline) ProcessEvaluation(ctx context.Context, cvPath, reportPath string) (*EvaluationResult, error) {
	log.Printf("Starting AI pipeline for CV: %s, Report: %s", cvPath, reportPath)

	// Step 1: Read and normalize file contents
	cvContent, err := p.fileReader.ReadFile(cvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CV file: %w", err)
	}

	reportContent, err := p.fileReader.ReadFile(reportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}

	log.Printf("Successfully read files - CV: %d chars, Report: %d chars", len(cvContent), len(reportContent))

	// Step 2: Stage 1 Analysis with Gemini
	stage1Analysis, err := p.geminiClient.Stage1Analysis(ctx, cvContent, reportContent)
	if err != nil {
		return nil, fmt.Errorf("failed Stage 1 analysis: %w", err)
	}

	log.Printf("Stage 1 analysis completed")

	// Step 3: Query ChromaDB for relevant context (optional)
	var chromaContext []string
	if p.chromaClient != nil {
		// Use a simple, relevant query for evaluation guidelines
		queryText := "CV evaluation guidelines project assessment scoring rubric"

		documents, err := p.chromaClient.QueryDocuments(ctx, queryText, 3)
		if err != nil {
			log.Printf("ChromaDB query failed, continuing without context: %v", err)
		} else {
			// Convert documents to string array
			for _, doc := range documents {
				chromaContext = append(chromaContext, doc.Content)
			}
			log.Printf("Retrieved %d context documents from ChromaDB", len(chromaContext))
		}
	}

	if len(chromaContext) == 0 {
		log.Printf("No ChromaDB context available, using default evaluation guidelines")
		// Provide default evaluation context when ChromaDB is not available
		chromaContext = []string{
			"CV Evaluation: Assess technical skills, experience level, education, and presentation quality. Rate CV match from 0.0-1.0.",
			"Project Evaluation: Assess code quality, complexity, documentation, and problem-solving approach. Rate project from 0.0-10.0.",
		}
	}

	// Step 4: Stage 2 Evaluation with context
	stage2Result, err := p.geminiClient.Stage2Evaluation(ctx, stage1Analysis, chromaContext, cvContent, reportContent)
	if err != nil {
		return nil, fmt.Errorf("failed Stage 2 evaluation: %w", err)
	}

	log.Printf("Stage 2 evaluation completed")

	// Step 5: Parse and return structured result
	result, err := p.parseEvaluationResult(stage2Result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse evaluation result: %w", err)
	}

	log.Printf("AI pipeline completed successfully")
	return result, nil
}

// parseEvaluationResult extracts JSON from Gemini response
func (p *Pipeline) parseEvaluationResult(response string) (*EvaluationResult, error) {
	// Find JSON content in the response
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[start : end+1]

	var result EvaluationResult
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Validate the result
	if result.CVMatchRate < 0 || result.CVMatchRate > 1 {
		result.CVMatchRate = 0.5 // Default fallback
	}
	if result.ProjectScore < 0 || result.ProjectScore > 10 {
		result.ProjectScore = 5.0 // Default fallback
	}

	return &result, nil
}
