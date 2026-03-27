# Manager Demo Script

## Goal
Show a complete narrative of a multi-agent, tool-augmented, HITL-gated workflow in <10 minutes.

## Demo Commands

```bash
# initialize profile/config
cc-agent-orchestration bootstrap

# inspect available tools
cc-agent-orchestration tools

# run orchestration simulation
cc-agent-orchestration simulate \
  --goal "implement feature flag rollout with slack + rag context"

# inspect event stream
cc-agent-orchestration events

# execute template workflow graph
cc-agent-orchestration pipeline

# show review gate behavior
cc-agent-orchestration review --confidence 0.77 --task-id task-2

# emit manager-facing report snapshot
cc-agent-orchestration report
```

## Talking Points

- Isolated worktree per agent session keeps branches clean and parallelizable.
- Provider abstraction lets us swap Claude/Codex/Aider/Gemini without core rewrites.
- Handoff messages are explicit and mode-driven (`summary`, `bugfix`, `review`).
- RAG integration is modeled as tools, not hardcoded into prompt logic.
- HITL gating and telemetry are first-class design concerns.

## Optional API Demo

```bash
cc-agent-orchestration server --addr :7070

curl -X POST http://localhost:7070/v1/run \
  -H 'content-type: application/json' \
  -d '{"goal":"triage and fix CI regression","trace_id":"trace-http-1"}'

curl http://localhost:7070/v1/events
```
