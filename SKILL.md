---
name: linear-cli
description: Read, create, update, comment on Linear issues from the shell. Use whenever the user asks to triage, file, assign, status-check, comment on, or delegate work in Linear.
---

# Linear CLI — Agent Usage

Binary: `linear` (Go, single static binary). All commands talk to Linear's GraphQL API.

## Auth

```sh
export LINEAR_API_KEY=lin_api_...
```

Or, to keep the key out of your context: `psst --global LINEAR_API_KEY -- linear ...`

Every command attributes actions to the API key's owner. Don't ghostwrite as a human — sign comments with your agent name.

## Output: always `--json`

Every command supports `--json`. Use it. Parse with `jq`. Never scrape the table output — column order is not stable across versions.

```sh
linear list --team ENG --json | jq '.[] | {id: .identifier, title}'
```

Other output modes (rarely needed): `--quiet` (id + url only), `--format id-only`, `--format tsv`, `--fields id,title,state.name` (TSV with selected fields).

## Discover identifiers once per session

Most write commands take a team key (`ENG`) or a project name. Look them up once and cache:

```sh
linear teams --json       # -> .[] | {key, id, name}
linear projects --json    # -> .[] | {id, name, slugId, status: .status.name}
linear states --team ENG --json    # workflow statuses for that team
linear labels --team ENG --json
linear users --json       # workspace members for assignment
linear me --json          # the API-key owner
```

`--team` wants the short key (`ENG`), not the full name.
`--project` accepts either a name (fuzzy, case-insensitive) or a UUID. Prefer UUID if you have it — names can collide.

## Read

```sh
linear list --team ENG --json -n 50
linear list --project "Q3 Launch" --status "In Progress" --json
linear list --assignee me --json
linear list -s "rate limit" --json            # full-text search (subset of fields)
linear search "auth bug" --json -n 20         # broader search, separate command
linear get ENG-142 --json                     # full issue + description + up to 100 comments
```

**`get` includes comments inline.** Use `linear get` instead of `linear list` + per-id `linear comments` when you need thread context.

`list` supports filters: `--team`, `--project`, `--status`, `--assignee` (`me` resolves to viewer), `-s/--search`, `-n/--limit` (default 20, max ~250).

## Write

```sh
# Create
linear create --team ENG --title "Fix login redirect" \
              --desc "Markdown body…" \
              --project "Q3 Launch" \
              --assignee "Alice" \
              --priority 2 --json

# Update — only the flags you pass get changed
linear update ENG-142 --status "In Review" --json
linear update ENG-142 --assignee me
linear update ENG-142 --project "Backlog"          # move between projects
linear update ENG-142 --clear-project              # detach from project
linear update ENG-142 --priority -1                # clear priority
linear update ENG-142 --due 2026-06-15T00:00:00Z   # ISO 8601
linear update ENG-142 --clear-due

# Comment
linear comment ENG-142 "Triaged, picked up by agent-koala"
linear comments ENG-142 --json                     # list existing

# Lifecycle
linear archive ENG-142
linear unarchive ENG-142
linear delete ENG-142                              # trash, recoverable from Linear UI
```

Priority: `1=urgent`, `2=high`, `3=medium`, `4=low`, `-1=clear`.
Status names are fuzzy-matched within the issue's team workflow — `"in prog"` will resolve to `"In Progress"` if unambiguous.
`--assignee me` resolves to the API key's owner.

## Batch create from JSON lines

```sh
cat <<'EOF' | linear batch-create --team ENG --project "Q3 Launch" --json
{"title": "Set up CI",   "priority": 2}
{"title": "Write tests", "priority": 3, "assignee": "Alice"}
{"title": "Cross-cutting task", "project": "Platform V3"}
EOF
```

Per-line `project` / `assignee` override the `--project` default. Per-line keys: `title` (required), `description`, `priority`, `project`, `assignee`.

## Common idioms

```sh
# Highest-priority unassigned ticket in a project
linear list --project "Q3 Launch" --json \
  | jq -r '[.[] | select(.assignee == null)] | sort_by(.priority) | .[0].identifier'

# Comment on every issue matching a search
linear list -s "flaky test" --json | jq -r '.[].identifier' \
  | while read id; do
      linear comment "$id" "agent-koala: investigating, tracked in #incident-42"
    done

# Bulk-close a sprint
linear list --project "Sprint 42" --status "In Review" --json \
  | jq -r '.[].identifier' \
  | xargs -I{} linear update {} --status "Done"

# Triage: pull title + last comment for everything assigned to me
for id in $(linear list --assignee me --status "Todo" --json | jq -r '.[].identifier'); do
  linear get "$id" --json | jq '{id: .identifier, title, last_comment: (.comments.nodes | last | .body)}'
done
```

## Command catalogue (reference)

Listed for discovery — most agent work needs only the bold ones.

| Category      | Commands                                                              |
| ------------- | --------------------------------------------------------------------- |
| **Issues**    | **list**, **get**, **create**, **update**, **comment**, **comments**, **search**, **batch-create**, archive, unarchive, delete |
| **Projects**  | **projects**, project-create, project-update                          |
| **Teams**     | **teams**                                                             |
| **Users**     | **users**, **me**                                                     |
| **States**    | **states**, state-create                                              |
| **Labels**    | labels, label-create, label-delete                                    |
| Cycles        | cycles, cycle-create                                                  |
| Initiatives   | initiatives, init-create                                              |
| Documents     | docs, doc-create                                                      |
| Webhooks      | webhooks, webhook-create, webhook-delete                              |
| Notifications | notifications, notif-archive, notif-read                              |

Run `linear <cmd> --help` for full flag list of any command.

## Failure modes worth handling

- `LINEAR_API_KEY not set` → don't retry, surface to the user.
- `team "FOO" not found` → call `linear teams` and present valid keys.
- `project "FOO" matched multiple: A, B, C` → ask the user which, or pass the UUID.
- `status "X" not found. Available: Todo, In Progress, ...` → workflow is team-specific; list `linear states --team ENG`.
- `user "X" not found` → call `linear users`; assignment uses exact name (or `me`).
- HTTP 429 / rate limit → back off, don't hammer.

## Don't

- Don't write `LINEAR_API_KEY` to files, logs, or commit messages.
- Don't parse human/table output — always `--json`.
- Don't loop with shrinking `--limit` to paginate; pass `-n` once (max ~250). For deeper scans, scope by team/project/status.
- Don't `delete` to "fix" a mistakenly-created issue while iterating — prefer `archive`, since delete is harder to recover for human reviewers.
- Don't impersonate a person in comments. Sign as your agent identity so humans can see what's automated.
