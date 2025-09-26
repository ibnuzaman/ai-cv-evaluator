package main

import (
	"aicvevaluator/internal/ai"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	ctx := context.Background()

	// Test Gemini client initialization
	client, err := ai.NewGeminiClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer client.Close()

	// Test a simple API call
	testCV := "Software Engineer with 3 years experience in Go, Python, and React."
	testReport := "Built a REST API using Go with PostgreSQL database and Docker deployment."

	result, err := client.Stage1Analysis(ctx, testCV, testReport)
	if err != nil {
		log.Fatalf("Failed to test Gemini API: %v", err)
	}

	fmt.Println("✅ Gemini API test successful!")
	fmt.Printf("Response length: %d characters\n", len(result))
	fmt.Println("Sample response (first 200 chars):")
	if len(result) > 200 {
		fmt.Printf("%s...\n", result[:200])
	} else {
		fmt.Println(result)
	}
}
