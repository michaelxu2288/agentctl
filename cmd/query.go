package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/michaelxu2288/cc-agent-orchestraiton/internal/pinecone"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query Pinecone to retrieve RAG context for agent workflows",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		namespace, _ := cmd.Flags().GetString("namespace")
		topK, _ := cmd.Flags().GetInt("top-k")
		vectorStr, _ := cmd.Flags().GetString("vector")

		apiKey := os.Getenv("PINECONE_API_KEY")
		client, err := pinecone.NewClient(apiKey, host)
		if err != nil {
			return err
		}

		vector, err := parseVector(vectorStr)
		if err != nil {
			return err
		}

		res, err := client.Query(context.Background(), pinecone.QueryRequest{
			Namespace: namespace,
			Vector:    vector,
			TopK:      topK,
		})
		if err != nil {
			return err
		}

		if len(res.Matches) == 0 {
			fmt.Println("no matches")
			return nil
		}

		for _, m := range res.Matches {
			fmt.Printf("%s\t%.6f\n", m.ID, m.Score)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().String("host", "", "Pinecone index host (e.g. https://<index>-<proj>.svc.<region>.pinecone.io)")
	queryCmd.Flags().String("namespace", "", "Pinecone namespace")
	queryCmd.Flags().Int("top-k", 5, "Top K matches")
	queryCmd.Flags().String("vector", "", "Comma-separated embedding vector (e.g. 0.1,0.2,0.3)")
	_ = queryCmd.MarkFlagRequired("host")
	_ = queryCmd.MarkFlagRequired("vector")
}

func parseVector(raw string) ([]float64, error) {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	vector := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid vector value %q: %w", p, err)
		}
		vector = append(vector, v)
	}
	if len(vector) == 0 {
		return nil, fmt.Errorf("vector is empty")
	}
	return vector, nil
}
