package chromadb

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	collectionName = "evaluation_guidelines"
	tenantID       = "default_tenant"
	databaseID     = "default_database"
)

// Client represents a ChromaDB client
type Client struct {
	baseURL      string
	httpClient   *http.Client
	collectionID string // Store the UUID of the collection
}

// Document represents a document in ChromaDB
type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"document"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Collection represents a ChromaDB collection
type Collection struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

// NewClient creates a new ChromaDB client
func NewClient(baseURL string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("ChromaDB URL is required")
	}

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return client, nil
}

// InitializeCollection initializes the evaluation guidelines collection using v2 API
func (c *Client) InitializeCollection(ctx context.Context) error {
	endpoint := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections", tenantID, databaseID)

	// Check if collection exists and get its UUID
	collection, err := c.getCollection(ctx, endpoint)
	if err == nil {
		c.collectionID = collection.ID
		fmt.Printf("✅ Collection '%s' already exists with ID: %s\n", collectionName, c.collectionID)
		return nil
	}

	// Create collection and get its UUID
	collection, err = c.createCollection(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	c.collectionID = collection.ID
	fmt.Printf("✅ Collection '%s' created successfully with ID: %s\n", collectionName, c.collectionID)
	return nil
}

// getCollection gets collection by name and returns its details including UUID
func (c *Client) getCollection(ctx context.Context, endpoint string) (*Collection, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to list collections: HTTP %d", resp.StatusCode)
	}

	var collections []Collection
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("failed to parse collections response: %w", err)
	}

	// Find our collection by name
	for _, collection := range collections {
		if collection.Name == collectionName {
			return &collection, nil
		}
	}

	return nil, fmt.Errorf("collection not found")
}

// createCollection creates collection and returns its details including UUID
func (c *Client) createCollection(ctx context.Context, endpoint string) (*Collection, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	reqBody := map[string]interface{}{
		"name": collectionName,
		"metadata": map[string]interface{}{
			"description": "Evaluation guidelines for CV and project evaluation",
		},
		"get_or_create": true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(respBody))
	}

	var collection Collection
	if err := json.NewDecoder(resp.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("failed to parse collection response: %w", err)
	}

	return &collection, nil
}

// generateSimpleEmbedding generates a simple embedding vector for text
func (c *Client) generateSimpleEmbedding(text string) []float64 {
	// Create a simple embedding based on text characteristics
	embedding := make([]float64, 384) // Common embedding size

	// Hash-based approach for consistent embeddings
	hash := md5.Sum([]byte(strings.ToLower(text)))

	// Convert hash bytes to float values
	for i := 0; i < len(embedding); i++ {
		byteIndex := i % len(hash)
		// Normalize to [-1, 1] range
		embedding[i] = (float64(hash[byteIndex]) - 127.5) / 127.5
	}

	// Add some text-based features
	textLen := float64(len(text))
	wordCount := float64(len(strings.Fields(text)))

	if len(embedding) > 2 {
		embedding[0] = textLen / 1000.0  // Normalize text length
		embedding[1] = wordCount / 100.0 // Normalize word count
	}

	return embedding
}

// AddDocument adds a document to the ChromaDB collection using collection UUID
func (c *Client) AddDocument(ctx context.Context, id, content string, metadata map[string]interface{}) error {
	if c.collectionID == "" {
		return fmt.Errorf("collection not initialized - no collection ID")
	}

	endpoint := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections/%s/add", tenantID, databaseID, c.collectionID)
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Generate embedding for the content
	embedding := c.generateSimpleEmbedding(content)

	reqBody := map[string]interface{}{
		"ids":        []string{id},
		"documents":  []string{content},
		"metadatas":  []map[string]interface{}{metadata},
		"embeddings": [][]float64{embedding},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// QueryDocuments queries documents from ChromaDB using collection UUID
func (c *Client) QueryDocuments(ctx context.Context, queryText string, n int) ([]Document, error) {
	if c.collectionID == "" {
		return nil, fmt.Errorf("collection not initialized - no collection ID")
	}

	endpoint := fmt.Sprintf("/api/v2/tenants/%s/databases/%s/collections/%s/query", tenantID, databaseID, c.collectionID)
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Generate embedding for the query
	queryEmbedding := c.generateSimpleEmbedding(queryText)

	reqBody := map[string]interface{}{
		"query_embeddings": [][]float64{queryEmbedding},
		"n_results":        n,
		"include":          []string{"documents", "metadatas", "distances"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		IDs       [][]string                 `json:"ids"`
		Documents [][]string                 `json:"documents"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
		Distances [][]float64                `json:"distances"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to Document slice
	documents := make([]Document, 0)
	if len(result.Documents) > 0 && len(result.Documents[0]) > 0 {
		for i := range result.Documents[0] {
			doc := Document{
				ID:      result.IDs[0][i],
				Content: result.Documents[0][i],
			}
			if len(result.Metadatas) > 0 && len(result.Metadatas[0]) > i {
				doc.Metadata = result.Metadatas[0][i]
			}
			documents = append(documents, doc)
		}
	}

	return documents, nil
}
