---
description: "Start or continue a Ratatoskr task: /task T-###"
argument-hint: "T-###"
---

Work on task $ARGUMENTS.

1. Open `docs/tasks/$ARGUMENTS-*.md`. If it does not exist, stop and say so.
2. Read AGENTS.md (already loaded) and only the spec sections listed in the task's `spec` field.
3. Check `depends_on` tasks are merged on main (`git log main --oneline | grep T-0xx` or their status is done). If not, stop.
4. If status is `todo`: set `status: in_progress`, `owner` (your tool/model), write the Plan section, commit `chore($ARGUMENTS): start`.
   If `changes_requested`: address every unchecked Review finding, tick them, reply inline if you disagree.
5. Implement within the task's `paths`, small Conventional Commits scoped `($ARGUMENTS)`.
6. Run `make check` until green. Fill Handoff. Set `status: in_review`. Summarize what you did in ≤10 lines.
