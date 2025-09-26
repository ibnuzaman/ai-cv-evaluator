package main

import (
	"context"
	"log"

	database "aicvevaluator/database/migration"
	"aicvevaluator/internal/ai"
	"aicvevaluator/internal/chromadb"
	"aicvevaluator/internal/config"
	"aicvevaluator/internal/handler"
	"aicvevaluator/internal/repository"
	"aicvevaluator/internal/service"
	"aicvevaluator/internal/util"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Initialize Database Connection with Migration
	ctx := context.Background()
	database.InitPostgresql(ctx, cfg)

	// 3. Get Database Connection
	db := database.GetPostgresql()
	if db == nil {
		log.Fatalf("failed to get database connection")
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// 3. Initialize AI Components
	fileReader := util.NewFileReader()

	// Initialize ChromaDB
	chromaClient, err := chromadb.NewClient(cfg.ChromaDBURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to ChromaDB: %v", err)
		log.Printf("Continuing without ChromaDB support...")
		chromaClient = nil
	} else {
		err = chromaClient.InitializeCollection(ctx)
		if err != nil {
			log.Printf("Warning: Failed to initialize ChromaDB collection: %v", err)
		}
	}

	// Initialize Gemini AI
	geminiClient, err := ai.NewGeminiClient(ctx, cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}
	defer geminiClient.Close()

	// Initialize AI Pipeline
	aiPipeline := ai.NewPipeline(fileReader, chromaClient, geminiClient)

	// 4. Initialize Layers (Dependency Injection)
	evaluationRepo := repository.NewEvaluationRepository(db)
	evaluationService := service.NewEvaluationService(evaluationRepo, aiPipeline)
	evaluationHandler := handler.NewEvaluationHandler(evaluationService)

	// 5. Setup Fiber App and Routes
	app := fiber.New()

	api := app.Group("/api/v1") // Grouping routes
	api.Post("/evaluate", evaluationHandler.Evaluate)
	api.Get("/result/:id", evaluationHandler.GetResult)
	// TODO: Add /upload endpoint later

	// 6. Start Server
	log.Printf("Server starting on port %s", cfg.AppPort)
	err = app.Listen(cfg.AppPort)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
