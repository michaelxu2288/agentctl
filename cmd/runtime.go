package cmd

import (
	"github.com/michaelxu2288/agentctl/internal/orchestration"
	"github.com/michaelxu2288/agentctl/internal/providers"
)

func newProviderRegistry() *providers.Registry {
	r := providers.NewRegistry()
	r.Register(providers.ClaudeAdapter{})
	r.Register(providers.CodexAdapter{})
	r.Register(providers.AiderAdapter{})
	r.Register(providers.GeminiAdapter{})
	return r
}

func newEngine() *orchestration.Engine {
	registry := newProviderRegistry()
	engine := orchestration.NewEngine(registry)
	engine.Router = orchestration.Router{
		PlannerSession:  "planner-claude",
		CoderSession:    "coder-codex",
		ReviewerSession: "reviewer-claude",
	}
	return engine
}
