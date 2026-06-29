# agentctl

> Terminal-native control plane for an autonomous multi-agent coding swarm.

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Build](https://img.shields.io/badge/build-passing-2ea44f)](#)
[![Tests](https://img.shields.io/badge/tests-passing-2ea44f)](#)
[![TUI](https://img.shields.io/badge/TUI-Bubble%20Tea-ff69b4)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-POC-orange)](#status)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

`agentctl` schedules a **swarm** of coding agents (Claude Code / Codex / Aider / Gemini)
across isolated git worktrees and PTY/tmux sessions. A **master coordinator** dispatches
work to **worker agents** over an internal gossip bus, agents hand off context to each
other, and the whole fleet renders as a live **kanban board** in your terminal. Buzzword
edition, but it builds, tests, and runs.

```
                          ┌────────────────────────────────────────┐
                          │              MASTER (coordinator)        │
                          │  scheduler · router · failure detector   │
                          └───────┬───────────────┬───────────────┬──┘
                  assign │ gossip │       assign  │      assign   │   internal bus (pub/sub)
                          ▼               ▼               ▼
                   ┌───────────┐   ┌───────────┐   ┌───────────┐
                   │ WORKER 1  │   │ WORKER 2  │   │ WORKER 3  │   ← each = git worktree + PTY/tmux
                   │ planner   │   │ coder     │   │ reviewer  │     + provider adapter
                   └─────┬─────┘   └─────┬─────┘   └─────┬─────┘
                         └── handoff ────┴── heartbeat ──┘
                                         │
                  BACKLOG │ RUNNING │ REVIEW │ DONE   ← kanban TUI
```

## Highlights

- **Master/worker control plane** — single coordinator, N workers, least-loaded dispatch with capability filtering, SWIM-style heartbeat failure detection + automatic requeue.
- **Internal agent comms** — in-process pub/sub bus with direct + broadcast routing, replayable history; clean seam for NATS/Redis transport.
- **Kanban TUI** — Bubble Tea board, cards flow backlog → running → review → done, vim keys (`h/j/k/l`, space to advance).
- **PTY/tmux sessions** — one attachable terminal per agent, pane capture for telemetry.
- **Git worktree isolation** — every agent gets its own branch + worktree; no stepping on each other.
- **Prompt handoff hooks** — `summary` / `bugfix` / `review` transformations pipe one agent's output into the next.
- **Pluggable providers** — Claude, Codex, Aider, Gemini adapters behind one interface with role/tool policies.
- **RAG tools** — LangGraph-style retrieval wrapper over a Pinecone HTTP client, callable mid-workflow.
- **HTTP control plane** — `/healthz`, `/v1/run`, `/v1/events`; metrics, traces, audit timeline.

## Tech stack

Go · Cobra · Bubble Tea / Lip Gloss · tmux · git worktrees · Pinecone · LangGraph-style RAG · Slack MCP

## Status

POC. The orchestration, swarm, kanban, and PTY layers compile and are unit-tested; provider/RAG/MCP adapters are wired against real interfaces but stubbed for demo. `tmux` + `git` required for live session commands.

## Quick start

```bash
go build -o agentctl .

# fan a goal across a master + 3 workers, tail the bus
./agentctl swarm --workers 3 --goal "ship the feature"

# open the kanban board
./agentctl kanban

# launch an agent in an isolated worktree + tmux
./agentctl orchestrate --name planner --provider claude --program "claude"

# hand off context planner -> coder
./agentctl handoff --source planner --target coder --mode bugfix

# inspect live PTY panes
./agentctl pty ls
./agentctl pty capture planner --lines 40
```

## Architecture

- `internal/swarm/` — master coordinator, worker, registry, internal bus (NEW)
- `internal/tui/kanban/` — Bubble Tea kanban board (NEW)
- `internal/pty/` — PTY/tmux session backend (NEW)
- `internal/orchestration/` — planner/router/scheduler/supervisor, event bus, retry/state
- `internal/providers/` — Claude/Codex/Aider/Gemini adapters + policy
- `internal/agent/` · `internal/worktree/` — session lifecycle + git worktree provisioning
- `internal/workflow/` · `internal/hooks/` — graph executor, templates, handoff transforms
- `internal/rag/` · `internal/pinecone/` · `internal/integrations/` — retrieval + Slack MCP
- `internal/api/` · `internal/telemetry/` — HTTP control plane, metrics/traces/audit

## Tests

```bash
go test ./...   # swarm, kanban, orchestration, providers, workflow, review, tools, memory
```

MIT. Built as a portfolio POC for distributed agent orchestration patterns in Go.
