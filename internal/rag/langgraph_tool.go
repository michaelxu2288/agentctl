package rag

import (
	"context"
	"fmt"

	"github.com/michaelxu2288/swarmboard/internal/pinecone"
)

type LangGraphRetriever struct {
	Pinecone *pinecone.Client
}

type RetrievalResult struct {
	IDs []string
}

// RetrieveContext simulates a LangGraph-style tool node that fetches top-k IDs from Pinecone.
// This keeps the implementation in Go while preserving the agent-tool contract used in workflows.
func (r *LangGraphRetriever) RetrieveContext(ctx context.Context, vector []float64, namespace string, topK int) (*RetrievalResult, error) {
	if r.Pinecone == nil {
		return nil, fmt.Errorf("pinecone client is nil")
	}

	res, err := r.Pinecone.Query(ctx, pinecone.QueryRequest{
		Namespace: namespace,
		Vector:    vector,
		TopK:      topK,
	})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(res.Matches))
	for _, m := range res.Matches {
		ids = append(ids, m.ID)
	}

	return &RetrievalResult{IDs: ids}, nil
}
