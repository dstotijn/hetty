# 07 — Multi-Agent Workflow

Several AI agents (Claude Code and OpenCode, different models) work on this repo **in parallel**. This document is
the protocol that keeps them from colliding and keeps quality consistent. A human ("operator") assigns tasks and merges.

## 1. Roles

| Role | Agent file | Writes | Must not | Suggested model class |
|---|---|---|---|---|
| **lead** (architect) | `lead.md` | `docs/**`, ADRs, task files, GraphQL schema (`internal/api/graphql/schema/**`), `go.mod`/`package.json` deps, CI | implement features outside contract/scaffold tasks | strongest reasoning model (e.g. Claude Opus in Claude Code) |
| **backend** | `backend.md` | `cmd/**`, `internal/**` except schema files, Go tests | change schema, touch `web/**` | strong Go coding model |
| **frontend** | `frontend.md` | `web/**` except `web/src/lib/api/generated/**` | change schema, touch `internal/**` | strong TS/Svelte coding model |
| **reviewer** | `reviewer.md` | review notes in the task file only | edit code | **different model family from the implementer** |
| **qa** | `qa.md` | tests only (`*_test.go`, `testdata/**`, `web/**/*.test.ts`, `web/tests/e2e/**`) | change production code | any solid model; can run in parallel with implementation |

Cross-model review is mandatory: if Claude implemented it, a non-Claude model reviews it in OpenCode, and vice versa.

## 2. Single source of truth
- `AGENTS.md` (root) — read automatically by OpenCode and by Claude Code (via `CLAUDE.md` → `@AGENTS.md` import).
- Specs in `docs/`. **If code and spec disagree, stop and report** — the lead updates the spec first.
- Each agent session loads only: `AGENTS.md` + its task file + the spec sections the task links. Don't read all docs.

## 3. Task files
- One file per task: `docs/tasks/T-###-slug.md`, created by the lead from `docs/tasks/TEMPLATE.md`.
- Frontmatter: `id, title, milestone, role, status, depends_on, owner, branch, paths, spec`.
- `status` lifecycle: `todo → in_progress → in_review → changes_requested → in_review → done` (or `blocked`).
- `paths` = the files/dirs the task may modify (**ownership**). Touching anything else requires a note in the task file
  ("Out-of-scope change needed: …") and lead approval — do not just do it.
- `docs/08-roadmap.md` holds the index table; the task file is authoritative for status.

## 4. Parallel work with git worktrees
```
# operator, once per task:
git worktree add ../ratatoskr-T-014 -b task/T-014-httpql-parser main
cd ../ratatoskr-T-014 && claude        # or: opencode
```
- One agent ↔ one worktree ↔ one task ↔ one branch. Never two agents in the same worktree.
- Rebase on `main` before review: `git fetch && git rebase origin/main`. Generated files conflicting? Re-run `make gen`, never hand-merge.
- Shared data dirs: use `--data-dir ./.dev-data` inside the worktree and a random port (`--listen 127.0.0.1:0`) so parallel runs don't clash.

## 5. Contract-first sequencing
1. Lead lands **contract tasks** first: GraphQL schema additions + stub resolvers returning `ErrNotImplemented` + generated types (both sides) + mock data for the UI.
2. Then backend and frontend tasks for that feature run **in parallel** against the contract; frontend uses MSW-style mocks (`web/src/lib/api/mocks`) until the backend lands.
3. QA writes tests from the acceptance criteria in parallel (tests may be committed as `t.Skip("pending T-0xx")` until the implementation merges).
4. Integration task (small) wires both and adds the E2E spec.

## 6. Task execution protocol (every implementer)
1. Read `AGENTS.md`, the task file, linked spec sections. Check `depends_on` are `done` (on `main`).
2. Set `status: in_progress`, `owner: <tool>/<model>`, commit `chore(T-###): start`.
3. Write a short **plan** in the task file (≤ 15 lines) before coding. If the plan reveals a spec gap → `status: blocked` + question.
4. Implement in small commits (Conventional Commits, scope = task id): `feat(T-014): parse clause operators`.
5. Run `make check` (and `make test-int` when touching proxy/storage). All green, no skipped lint.
6. Fill **Handoff** section: what changed, how to verify, known limitations, follow-ups (as proposed new task stubs).
7. Set `status: in_review`. Push the branch.

## 7. Review protocol (reviewer)
- Run `/review T-###` (command file provided). Review the diff `main...task/T-###` against: acceptance criteria, `AGENTS.md` hard rules,
  conventions docs, tests quality, security (path traversal, injection into SQL, unbounded memory, goroutine leaks), performance budgets.
- Write findings into the task file under **Review** as a checklist with severity `blocker|major|minor|nit` and file:line.
- Verdict: `approve` or `changes_requested`. Max 3 review rounds; then escalate to lead.

## 8. Pull request & merge (operator / lead)
- When a task reaches `in_review`, the implementer pushes the branch and opens a PR **on our fork only**:
  `gh pr create --repo OWNER/ratatoskr --base main --title "T-014: <title>" --body-file <(sed -n '/## Handoff/,$p' docs/tasks/T-014-*.md)`.
- CI must be green. The reviewer agent writes findings in the task file (source of truth) and may mirror the summary as a PR comment
  (`gh pr comment --repo OWNER/ratatoskr`).
- Operator squash-merges in dependency order (`gh pr merge --repo OWNER/ratatoskr --squash --delete-branch`), sets task `status: done`
  on `main`, removes the worktree: `git worktree remove ../ratatoskr-T-014`. Update roadmap table.

## 9. Context hygiene & cost control
- Prefer many small tasks (½–1 day of human-equivalent work, ≤ ~600 changed lines).
- Long sessions drift: if a session passes ~60% context, write the Handoff section and start a fresh session.
- Agents must not paste large generated files into chat; reference paths.
- Use cheaper models for QA/test writing and nits; strongest models for lead/review of core packages (`rawhttp`, `proxy`, `query`, `store`, `workflow`).

## 10. Definition of Done (all tasks)
- [ ] Acceptance criteria met and demonstrated by tests.
- [ ] `make check` green; integration/E2E added when the task touches user-visible behaviour.
- [ ] No TODO without a task id (`// TODO(T-123): …`).
- [ ] Docs updated if behaviour/API changed (spec, schema comments, user-facing help text).
- [ ] Handoff section complete; reviewer approved.
