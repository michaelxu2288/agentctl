# Multi Agent Orchestration Terminal App

Terminal-first orchestration framework for running concurrent Claude Code/Codex-style workflows with isolated git worktrees, tmux-backed agent sessions, prompt handoff hooks, and Pinecone-backed retrieval tools.

## Tech Stack

- Go
- Cobra
- Pinecone
- CLI/TUI

## Resume-Aligned Scope

- Built a TUI for orchestrating concurrent agentic workflows across git worktrees with PTY/tmux session management.
- Designed multi-agent launches for Claude Code/Codex CLI with prompt handoff and hook-based piping between chats.
- Integrated Slack MCP with LangGraph-style vector retrieval for agents to query RAG pipelines as callable tools mid-workflow.

## Architecture

- `cmd/`: Cobra command surface (`orchestrate`, `handoff`, `query`, `rag`, `mcp`, `tui`, etc.)
- `internal/agent/`: session model, tmux lifecycle, prompt routing, state persistence
- `internal/worktree/`: git worktree provisioning + branch lifecycle
- `internal/hooks/`: handoff hook prompt transformations (`summary`, `bugfix`, `review`)
- `internal/pinecone/`: Pinecone HTTP query client
- `internal/rag/`: LangGraph-style retrieval tool wrapper over Pinecone
- `internal/integrations/`: Slack MCP tool-call client
- `internal/tui/`: Bubble Tea session board for interactive monitoring

## Quick Start

```bash
# 1) Build
GO111MODULE=on go build -o cc-agent-orchestration .

# 2) Launch an agent in isolated branch/worktree
./cc-agent-orchestration orchestrate \
  --repo /path/to/repo \
  --name planner-claude \
  --provider claude \
  --program "claude"

# 3) Launch a second agent
./cc-agent-orchestration orchestrate \
  --repo /path/to/repo \
  --name coder-codex \
  --provider codex \
  --program "codex"

# 4) Pipe context from planner -> coder using hook mode
./cc-agent-orchestration handoff \
  --repo /path/to/repo \
  --source planner-claude \
  --target coder-codex \
  --mode bugfix

# 5) Open TUI board
./cc-agent-orchestration tui --repo /path/to/repo
```

## Command Reference

```bash
cc-agent-orchestration orchestrate --name <session> --program "claude|codex|..."
cc-agent-orchestration list
cc-agent-orchestration prompt --name <session> --text "<prompt>"
cc-agent-orchestration handoff --source <sessionA> --target <sessionB> --mode summary|bugfix|review
cc-agent-orchestration kill --name <session>
cc-agent-orchestration tui
```

## Pinecone + RAG Tooling

Set your Pinecone API key:

```bash
export PINECONE_API_KEY=<your_key>
```

Query Pinecone directly:

```bash
cc-agent-orchestration query \
  --host https://<index-host> \
  --namespace prod \
  --vector 0.12,0.34,0.56 \
  --top-k 5
```

Run LangGraph-style retrieval wrapper:

```bash
cc-agent-orchestration rag \
  --host https://<index-host> \
  --namespace prod \
  --vector 0.12,0.34,0.56 \
  --top-k 5
```

## Slack MCP Tool Calls

Set your MCP token:

```bash
export SLACK_MCP_TOKEN=<your_token>
```

Invoke an MCP tool:

```bash
cc-agent-orchestration mcp \
  --endpoint https://<mcp-gateway> \
  --tool slack.search_messages \
  --params channel=C123,query=deploy-failure
```

## Notes

- This repo intentionally emphasizes implementation patterns for agent orchestration and tool wiring in Go.
- `tmux` and `git` must be installed locally for session/worktree orchestration commands.
