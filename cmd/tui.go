package cmd

import (
	"github.com/michaelxu2288/swarmboard/internal/agent"
	"github.com/michaelxu2288/swarmboard/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open interactive session board for multi-agent orchestration",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr, err := agent.NewManager(repoPath)
		if err != nil {
			return err
		}
		sessions, err := mgr.List()
		if err != nil {
			return err
		}
		return tui.Run(sessions)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
