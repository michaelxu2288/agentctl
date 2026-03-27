package tools

import "testing"

func TestToolRegistry(t *testing.T) {
	r := NewRegistry()
	r.Register(PineconeQueryTool{})
	r.Register(SlackSearchTool{})
	if len(r.Names()) != 2 {
		t.Fatalf("expected 2 tools")
	}
	t1, err := r.Get("pinecone.query")
	if err != nil {
		t.Fatalf("expected pinecone tool")
	}
	if _, err := t1.Call(map[string]any{"namespace": "prod"}); err != nil {
		t.Fatalf("unexpected call error: %v", err)
	}
}
