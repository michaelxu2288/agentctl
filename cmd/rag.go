package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/michaelxu2288/cc-agent-orchestraiton/internal/pinecone"
	"github.com/michaelxu2288/cc-agent-orchestraiton/internal/rag"
	"github.com/spf13/cobra"
)

var ragCmd = &cobra.Command{
	Use:   "rag",
	Short: "Run LangGraph-style retrieval step backed by Pinecone",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		namespace, _ := cmd.Flags().GetString("namespace")
		topK, _ := cmd.Flags().GetInt("top-k")
		vectorStr, _ := cmd.Flags().GetString("vector")

		apiKey := os.Getenv("PINECONE_API_KEY")
		pc, err := pinecone.NewClient(apiKey, host)
		if err != nil {
			return err
		}

		vector, err := parseVector(vectorStr)
		if err != nil {
			return err
		}

		retriever := &rag.LangGraphRetriever{Pinecone: pc}
		result, err := retriever.RetrieveContext(context.Background(), vector, namespace, topK)
		if err != nil {
			return err
		}

		for _, id := range result.IDs {
			fmt.Println(id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ragCmd)
	ragCmd.Flags().String("host", "", "Pinecone index host")
	ragCmd.Flags().String("namespace", "", "Namespace")
	ragCmd.Flags().Int("top-k", 5, "Top K results")
	ragCmd.Flags().String("vector", "", "Comma-separated embedding vector")
	_ = ragCmd.MarkFlagRequired("host")
	_ = ragCmd.MarkFlagRequired("vector")
}
