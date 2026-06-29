package cmd

import (
	"fmt"

	"github.com/michaelxu2288/agentctl/internal/agent"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill an orchestrated agent session and clean up worktree",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		mgr, err := agent.NewManager(repoPath)
		if err != nil {
			return err
		}
		if err := mgr.Kill(name); err != nil {
			return err
		}
		fmt.Printf("killed session: %s\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().String("name", "", "Session name")
	_ = killCmd.MarkFlagRequired("name")
}
