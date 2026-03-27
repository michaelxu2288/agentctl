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
	Use:   "cc-agent-orchestration",
	Short: "Multi-agent orchestration terminal app for Claude Code-style workflows",
	Long:  "Orchestrate concurrent AI agent workflows with git worktrees, tmux sessions, prompt handoffs, and Pinecone retrieval tools.",
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
