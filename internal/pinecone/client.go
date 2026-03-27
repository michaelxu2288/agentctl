package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	APIKey string
	Host   string
}

type QueryRequest struct {
	Namespace string    `json:"namespace,omitempty"`
	Vector    []float64 `json:"vector"`
	TopK      int       `json:"topK"`
}

type Match struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}

type QueryResponse struct {
	Matches []Match `json:"matches"`
}

func NewClient(apiKey, host string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("pinecone api key is required")
	}
	if host == "" {
		return nil, fmt.Errorf("pinecone host is required")
	}
	return &Client{APIKey: apiKey, Host: strings.TrimRight(host, "/")}, nil
}

func (c *Client) Query(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	if len(req.Vector) == 0 {
		return nil, fmt.Errorf("vector must not be empty")
	}
	if req.TopK <= 0 {
		req.TopK = 5
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	httpClient := &http.Client{Timeout: 20 * time.Second}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Host+"/query", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Api-Key", c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("pinecone request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("pinecone error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out QueryResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("failed to parse pinecone response: %w", err)
	}
	return &out, nil
}
