# Manager Demo Script

## Goal
Show a complete narrative of a multi-agent, tool-augmented, HITL-gated workflow in <10 minutes.

## Demo Commands

```bash
# initialize profile/config
swarmboard bootstrap

# inspect available tools
swarmboard tools

# run orchestration simulation
swarmboard simulate \
  --goal "implement feature flag rollout with slack + rag context"

# inspect event stream
swarmboard events

# execute template workflow graph
swarmboard pipeline

# show review gate behavior
swarmboard review --confidence 0.77 --task-id task-2

# emit manager-facing report snapshot
swarmboard report
```

## Talking Points

- Isolated worktree per agent session keeps branches clean and parallelizable.
- Provider abstraction lets us swap Claude/Codex/Aider/Gemini without core rewrites.
- Handoff messages are explicit and mode-driven (`summary`, `bugfix`, `review`).
- RAG integration is modeled as tools, not hardcoded into prompt logic.
- HITL gating and telemetry are first-class design concerns.

## Optional API Demo

```bash
swarmboard server --addr :7070

curl -X POST http://localhost:7070/v1/run \
  -H 'content-type: application/json' \
  -d '{"goal":"triage and fix CI regression","trace_id":"trace-http-1"}'

curl http://localhost:7070/v1/events
```
