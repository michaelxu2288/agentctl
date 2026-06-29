package cmd

import (
	"fmt"

	"github.com/michaelxu2288/agentctl/internal/tools"
	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List available callable tools for agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := tools.NewRegistry()
		r.Register(tools.PineconeQueryTool{})
		r.Register(tools.SlackSearchTool{})
		r.Register(tools.CompositeTool{
			NameValue:        "context.bundle",
			DescriptionValue: "Bundle slack + rag context",
			Steps: []tools.Tool{
				tools.SlackSearchTool{},
				tools.PineconeQueryTool{},
			},
		})

		for _, name := range r.Names() {
			t, err := r.Get(name)
			if err != nil {
				continue
			}
			fmt.Printf("- %s: %s\n", t.Name(), t.Description())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}
