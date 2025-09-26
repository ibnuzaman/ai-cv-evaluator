package main

import (
	"aicvevaluator/internal/chromadb"
	"context"
	"log"
	"os"
)

func main() {
	chromaURL := os.Getenv("CHROMADB_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	client, err := chromadb.NewClient(chromaURL)
	if err != nil {
		log.Fatalf("Failed to create ChromaDB client: %v", err)
	}

	ctx := context.Background()
	err = client.InitializeCollection(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize collection: %v", err)
	}

	// Sample evaluation guidelines and context data
	contexts := []struct {
		ID       string
		Document string
		Metadata map[string]interface{}
	}{
		{
			ID: "backend_skills_golang",
			Document: `Go (Golang) Backend Development Skills:
- Strong understanding of Go syntax, goroutines, and channels
- Experience with popular Go frameworks like Gin, Echo, or Fiber
- Database integration with GORM, sqlx, or standard database/sql
- RESTful API design and implementation
- Microservices architecture knowledge
- Docker containerization and deployment
- Testing with Go testing package and testify
- Version control with Git
- Understanding of clean architecture and dependency injection`,
			Metadata: map[string]interface{}{
				"category": "backend_skills",
				"language": "golang",
				"level":    "intermediate",
			},
		},
		{
			ID: "project_evaluation_criteria",
			Document: `Project Evaluation Criteria:
1. Code Quality (25%):
   - Clean, readable, and well-structured code
   - Proper error handling and logging
   - Following Go best practices and conventions
   - Code documentation and comments

2. Architecture (25%):
   - Clean architecture implementation
   - Separation of concerns
   - Dependency injection
   - Database design and migrations

3. Functionality (25%):
   - Working REST API endpoints
   - CRUD operations implementation
   - Input validation and sanitization
   - Response formatting

4. Technical Implementation (25%):
   - Database integration and queries
   - Authentication and authorization (if required)
   - Testing coverage
   - Docker containerization`,
			Metadata: map[string]interface{}{
				"category": "evaluation_criteria",
				"type":     "project_assessment",
			},
		},
		{
			ID: "cv_evaluation_guidelines",
			Document: `CV Evaluation Guidelines:
1. Experience Level Assessment:
   - Junior (0-2 years): Basic understanding, simple projects
   - Mid-level (2-5 years): Solid experience, complex projects
   - Senior (5+ years): Leadership, architecture decisions, mentoring

2. Technical Skills Evaluation:
   - Programming languages proficiency
   - Framework and library experience
   - Database management skills
   - DevOps and deployment knowledge
   - Testing and quality assurance

3. Project Portfolio:
   - Diversity of projects
   - Complexity and scale
   - Technologies used
   - Problem-solving approach
   - Documentation quality

4. Soft Skills Indicators:
   - Communication skills from project descriptions
   - Teamwork and collaboration
   - Learning attitude and adaptability
   - Problem-solving mindset`,
			Metadata: map[string]interface{}{
				"category": "cv_evaluation",
				"type":     "guidelines",
			},
		},
		{
			ID: "scoring_rubric",
			Document: `Scoring Rubric:
CV Match Rate (0.0-1.0):
- 0.9-1.0: Exceptional match, exceeds requirements
- 0.8-0.9: Strong match, meets most requirements
- 0.6-0.8: Good match, meets basic requirements
- 0.4-0.6: Fair match, some gaps in requirements
- 0.2-0.4: Poor match, significant gaps
- 0.0-0.2: Very poor match, major misalignment

Project Score (0-10):
- 9-10: Exceptional quality, production-ready
- 7-8: High quality, minor improvements needed
- 5-6: Good quality, some improvements needed
- 3-4: Fair quality, significant improvements needed
- 1-2: Poor quality, major issues present
- 0: Very poor quality or non-functional`,
			Metadata: map[string]interface{}{
				"category": "scoring",
				"type":     "rubric",
			},
		},
	}

	for _, ctx := range contexts {
		err := client.AddDocument(context.Background(), ctx.ID, ctx.Document, ctx.Metadata)
		if err != nil {
			log.Printf("Warning: Failed to add document %s: %v", ctx.ID, err)
		} else {
			log.Printf("Added document: %s", ctx.ID)
		}
	}

	log.Println("ChromaDB seeding completed successfully!")
}
