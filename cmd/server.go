package cmd

import (
	"fmt"

	"github.com/michaelxu2288/cc-agent-orchestraiton/internal/api"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run HTTP API for orchestration control-plane",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		engine := newEngine()
		h := api.NewHandlers(engine)
		srv := api.NewServer(addr, h)
		fmt.Printf("api listening on %s\n", addr)
		return srv.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("addr", ":7070", "Bind address")
}
