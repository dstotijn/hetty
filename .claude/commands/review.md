---
description: "Review a Ratatoskr task branch: /review T-###"
argument-hint: "T-###"
---

Act as the reviewer agent for task $ARGUMENTS on the current branch.
Follow the reviewer procedure exactly (AGENTS.md + `.claude/agents/reviewer.md` / `.opencode/agents/reviewer.md`).
Write findings and verdict into the Review section of `docs/tasks/$ARGUMENTS-*.md` and commit it as `docs($ARGUMENTS): review`.
Do not modify any other file.
