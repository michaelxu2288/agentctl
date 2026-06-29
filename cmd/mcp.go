package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/michaelxu2288/swarmboard/internal/integrations"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Call Slack MCP tools from the orchestration CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, _ := cmd.Flags().GetString("endpoint")
		tool, _ := cmd.Flags().GetString("tool")
		paramsRaw, _ := cmd.Flags().GetString("params")

		token := os.Getenv("SLACK_MCP_TOKEN")
		client, err := integrations.NewSlackMCPClient(endpoint, token)
		if err != nil {
			return err
		}

		params := map[string]any{}
		if strings.TrimSpace(paramsRaw) != "" {
			for _, kv := range strings.Split(paramsRaw, ",") {
				parts := strings.SplitN(strings.TrimSpace(kv), "=", 2)
				if len(parts) != 2 {
					continue
				}
				params[parts[0]] = parts[1]
			}
		}

		res, err := client.CallTool(context.Background(), integrations.ToolCallRequest{
			Tool:   tool,
			Params: params,
		})
		if err != nil {
			return err
		}

		fmt.Println(res.Output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.Flags().String("endpoint", "", "Slack MCP gateway endpoint")
	mcpCmd.Flags().String("tool", "", "Tool name to call")
	mcpCmd.Flags().String("params", "", "Comma-separated k=v params (e.g. channel=C123,query=latest incident)")
	_ = mcpCmd.MarkFlagRequired("endpoint")
	_ = mcpCmd.MarkFlagRequired("tool")
}
