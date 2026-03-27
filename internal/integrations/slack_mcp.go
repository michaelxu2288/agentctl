package integrations

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

type SlackMCPClient struct {
	Endpoint string
	Token    string
}

type ToolCallRequest struct {
	Tool   string         `json:"tool"`
	Params map[string]any `json:"params"`
}

type ToolCallResponse struct {
	Output string `json:"output"`
}

func NewSlackMCPClient(endpoint, token string) (*SlackMCPClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("mcp endpoint is required")
	}
	if token == "" {
		return nil, fmt.Errorf("mcp token is required")
	}
	return &SlackMCPClient{Endpoint: strings.TrimRight(endpoint, "/"), Token: token}, nil
}

func (c *SlackMCPClient) CallTool(ctx context.Context, req ToolCallRequest) (*ToolCallResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mcp request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint+"/tools/call", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to build mcp request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 20 * time.Second}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mcp request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mcp error (%d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var out ToolCallResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("failed to decode mcp response: %w", err)
	}
	return &out, nil
}
