package chromadb

import (
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

// Client wraps ChromaDB operations
type Client struct {
	client     *chromago.Client
	collection *chromago.Collection
}

// NewClient creates a new ChromaDB client
func NewClient(url string) (*Client, error) {
	client, err := chromago.NewClient(chromago.WithBasePath(url))
	if err != nil {
		return nil, fmt.Errorf("failed to create ChromaDB client: %w", err)
	}

	return &Client{client: client}, nil
}

// InitializeCollection creates or gets the evaluation collection
func (c *Client) InitializeCollection(ctx context.Context) error {
	collectionName := "cv_evaluation_context"

	// Try to get existing collection first
	collection, err := c.client.GetCollection(ctx, collectionName, nil)
	if err != nil {
		// Collection doesn't exist, create it
		collection, err = c.client.CreateCollection(ctx, collectionName, nil, true, nil, types.L2)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}

	c.collection = collection
	return nil
}

// QuerySimilar retrieves similar documents from ChromaDB
func (c *Client) QuerySimilar(ctx context.Context, queryText string, nResults int) ([]string, error) {
	if c.collection == nil {
		return nil, fmt.Errorf("collection not initialized")
	}

	queryResult, err := c.collection.Query(ctx, []string{queryText}, int32(nResults), nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query ChromaDB: %w", err)
	}

	var results []string
	if len(queryResult.Documents) > 0 {
		for _, doc := range queryResult.Documents[0] {
			results = append(results, doc)
		}
	}

	return results, nil
}

// AddDocument adds a document to the collection
func (c *Client) AddDocument(ctx context.Context, id, document string, metadata map[string]interface{}) error {
	if c.collection == nil {
		return fmt.Errorf("collection not initialized")
	}

	metadataList := []map[string]interface{}{metadata}
	_, err := c.collection.Add(ctx, nil, metadataList, []string{document}, []string{id})
	if err != nil {
		return fmt.Errorf("failed to add document: %w", err)
	}

	return nil
}
