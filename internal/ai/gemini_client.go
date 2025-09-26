package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient wraps Gemini API operations
type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Use gemini-1.5-pro which is the current available model
	model := client.GenerativeModel("gemini-1.5-pro")
	model.SetTemperature(0.1) // Lower temperature for more consistent results

	// Set safety settings to be more permissive for business evaluation content
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (g *GeminiClient) Close() error {
	return g.client.Close()
}

// Stage1Analysis performs initial analysis of CV and project report
func (g *GeminiClient) Stage1Analysis(ctx context.Context, cvContent, reportContent string) (string, error) {
	prompt := fmt.Sprintf(`
You are an expert CV and project evaluator. Analyze the provided CV and project report.

CV Content:
%s

Project Report Content:
%s

Please provide an initial analysis focusing on:
1. Key skills and experience from the CV
2. Project complexity and technical depth
3. Alignment between CV skills and project requirements
4. Initial impressions and areas that need deeper evaluation

Provide a structured analysis in JSON format with the following structure:
{
  "cv_skills": ["skill1", "skill2", ...],
  "cv_experience_level": "junior/mid/senior",
  "project_complexity": "low/medium/high",
  "project_technologies": ["tech1", "tech2", ...],
  "skill_alignment": "poor/fair/good/excellent",
  "areas_for_deeper_evaluation": ["area1", "area2", ...]
}
`, cvContent, reportContent)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		// Check if it's a quota exceeded error
		if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") {
			return "", fmt.Errorf("Gemini API quota exceeded. Please check your billing or wait for quota reset: %w", err)
		}
		return "", fmt.Errorf("failed to generate Stage 1 analysis: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in Gemini API response")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in Gemini API response")
	}

	// Check if the response was blocked due to safety filters
	if resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
		return "", fmt.Errorf("response blocked by Gemini safety filters")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}

// Stage2Evaluation performs refined evaluation using context from ChromaDB
func (g *GeminiClient) Stage2Evaluation(ctx context.Context, stage1Analysis string, chromaContext []string, cvContent, reportContent string) (string, error) {
	contextStr := strings.Join(chromaContext, "\n\n")

	prompt := fmt.Sprintf(`
You are an expert CV and project evaluator. Based on the initial analysis and additional context, provide a comprehensive evaluation.

Initial Analysis:
%s

Additional Context from Knowledge Base:
%s

CV Content:
%s

Project Report Content:
%s

Based on all this information, provide a comprehensive evaluation in the following JSON format:
{
  "cv_match_rate": 0.0-1.0,
  "cv_feedback": "detailed feedback on CV quality, strengths, and areas for improvement",
  "project_score": 0.0-10.0,
  "project_feedback": "detailed feedback on project quality, technical implementation, and documentation",
  "overall_summary": "comprehensive summary of the candidate's suitability and recommendations"
}

Scoring Guidelines:
- cv_match_rate: How well the CV matches the project requirements (0.0 = no match, 1.0 = perfect match)
- project_score: Overall project quality (0-10 scale, where 10 is exceptional)

Provide constructive, specific feedback that helps the candidate improve.
`, stage1Analysis, contextStr, cvContent, reportContent)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		// Check if it's a quota exceeded error
		if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") {
			return "", fmt.Errorf("Gemini API quota exceeded. Please check your billing or wait for quota reset: %w", err)
		}
		return "", fmt.Errorf("failed to generate Stage 2 evaluation: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in Gemini API response")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in Gemini API response")
	}

	// Check if the response was blocked due to safety filters
	if resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
		return "", fmt.Errorf("response blocked by Gemini safety filters")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
