package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	repoPath string
)

var rootCmd = &cobra.Command{
	Use:   "agentctl",
	Short: "Control plane for an autonomous multi-agent swarm",
	Long:  "agentctl is a terminal-native control plane for orchestrating swarms of coding agents across isolated git worktrees, PTY/tmux sessions, a master/worker scheduler, an internal gossip bus, and a kanban TUI.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", cwd, "Path to the git repository used for orchestration")
}
