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

	// Step 3: Query ChromaDB for relevant context
	queryText := fmt.Sprintf("%s %s", cvContent, reportContent)
	// Limit query text to avoid too long requests
	if len(queryText) > 1000 {
		queryText = queryText[:1000]
	}

	chromaContext, err := p.chromaClient.QuerySimilar(ctx, queryText, 5)
	if err != nil {
		log.Printf("ChromaDB query failed, continuing without context: %v", err)
		chromaContext = []string{} // Continue without context if ChromaDB fails
	}

	log.Printf("Retrieved %d context documents from ChromaDB", len(chromaContext))

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
