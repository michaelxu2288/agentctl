package tools

import "fmt"

type SlackSearchTool struct{}

func (s SlackSearchTool) Name() string { return "slack.search" }
func (s SlackSearchTool) Description() string {
	return "Search Slack context via MCP bridge for incidents and team notes"
}
func (s SlackSearchTool) Call(params map[string]any) (string, error) {
	query, _ := params["query"].(string)
	if query == "" {
		query = "latest incidents"
	}
	return fmt.Sprintf("mock slack search results for query=%q", query), nil
}
