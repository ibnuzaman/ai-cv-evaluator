package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"aicvevaluator/internal/chromadb"
	"aicvevaluator/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	dataDir := flag.String("dir", "data/evaluation_guidelines", "Directory containing evaluation guideline documents")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create ChromaDB client
	client, err := chromadb.NewClient(cfg.ChromaDBURL)
	if err != nil {
		log.Fatalf("Failed to create ChromaDB client: %v", err)
	}

	// Initialize collection (this will create it if it doesn't exist)
	ctx := context.Background()
	if err := client.InitializeCollection(ctx); err != nil {
		log.Fatalf("Failed to initialize ChromaDB collection: %v", err)
	}
	fmt.Println("✅ ChromaDB collection initialized")

	// Read files from the data directory
	files, err := os.ReadDir(*dataDir)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
	}

	// Process each file
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue // Skip directories and non-text files
		}

		filePath := filepath.Join(*dataDir, file.Name())
		fmt.Printf("Processing %s...\n", filePath)

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: Failed to read file %s: %v", filePath, err)
			continue
		}

		// Add document to ChromaDB
		docID := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		metadata := map[string]interface{}{
			"source": filePath,
			"type":   "evaluation_guideline",
		}

		err = client.AddDocument(ctx, docID, string(content), metadata)
		if err != nil {
			log.Printf("Warning: Failed to add document %s to ChromaDB: %v", docID, err)
			continue
		}

		fmt.Printf("✅ Added document %s to ChromaDB\n", docID)
	}

	fmt.Println("✅ Seeding completed successfully!")
}
