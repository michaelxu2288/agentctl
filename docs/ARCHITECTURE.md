# Architecture Overview

## Positioning

This MVP is intentionally implementation-heavy to demonstrate a credible multi-agent orchestration surface area:
- orchestration engine
- provider abstraction and adapters
- git worktree + tmux session lifecycle
- prompt handoff bus and hook policies
- HITL review gate
- Slack MCP + Pinecone/LangGraph-style retrieval tools
- API control plane and TUI observer

## Runtime Layers

1. CLI (`cmd/`)
- Cobra command tree for orchestration, pipelines, tools, reporting, and API service.

2. Orchestration (`internal/orchestration/`)
- Planning/decomposition
- Routing by role
- Scheduling by priority and concurrency
- Event bus + handoff queue
- Retry policies and state snapshots

3. Providers (`internal/providers/`)
- Provider interface with role/tool capability checks
- Concrete adapters for Claude/Codex/Aider/Gemini
- Policy helper for tool allow-lists

4. Workflow (`internal/workflow/`)
- Directed graph nodes and dependencies
- Template graph for Claude->Codex->Review->HITL pipeline
- Executor with deterministic ordering for MVP demos

5. Tools (`internal/tools/`)
- Tool registry pattern
- Pinecone query + Slack search tools
- Composite tool for bundled context retrieval

6. Integrations (`internal/integrations/`, `internal/pinecone/`, `internal/rag/`)
- Slack MCP HTTP client
- Pinecone HTTP query client
- LangGraph-style retrieval wrapper

7. API (`internal/api/`)
- Health, run, events endpoints
- Designed as control-plane entry point for external dashboards

8. Observability (`internal/telemetry/`)
- Counter metrics
- Trace sink (JSONL)
- Audit trail renderer

## Data Objects

- `Task`: queued/running/reviewing lifecycle
- `Event`: internal telemetry events and actor targeting
- `Artifact`: generated patches/notes/tool outputs
- `SessionSnapshot`: session routing and capability metadata
- `Snapshot`: persisted run state timeline

## Typical Flow

1. `simulate` command sends goal to engine.
2. Planner creates tasks.
3. Scheduler orders tasks.
4. Router assigns each task to planner/coder/reviewer sessions.
5. Provider adapter generates role-specific completion output.
6. Suggested-next payloads are converted into handoff queue messages.
7. Events and optional state snapshots are persisted for replay/reporting.

## Why This Shape

- Mirrors real production orchestration architecture while staying MVP-friendly.
- Makes it easy to replace mock provider responses with live SDK/CLI invocations.
- Keeps manager-facing demos simple (`simulate`, `pipeline`, `report`, `server`).
