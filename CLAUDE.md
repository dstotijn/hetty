@AGENTS.md

## Claude Code specifics
- Subagents for this repo live in `.claude/agents/` (lead, backend, frontend, reviewer, qa). Pick the one matching your task's `role`.
- Slash commands: `/task T-###` (start/continue a task), `/review T-###` (review a branch), `/handoff` (wrap up a session), `/plan-milestone M#` (lead only).
- Use plan mode for tasks touching `internal/rawhttp`, `internal/proxy`, `internal/query`, `internal/store`, `internal/workflow`.
- Keep `docs/tasks/*` updated as you go; it is the shared memory with OpenCode agents. Don't rely on your own auto-memory for project state.
- Graphify: `graphify-out/GRAPH_REPORT.md` and `graphify query "<question>"` give a quick map of the codebase. Use them before broad searches. `/graphify . --update` refreshes the graph. See `AGENTS.md` for details.
