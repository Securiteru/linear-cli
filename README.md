# linear-cli

Fast, native Go CLI for the Linear issue tracker via GraphQL. Zero dependencies beyond cobra -- single static binary, no runtime, no Node, no Python.

Built for agentic use: works with [psst](https://github.com/nicois/psst) for secret injection so API keys never leak into agent context.

## Features

**Issues** -- `list`, `get`, `create`, `update`, `delete`, `archive`, `unarchive`, `search`, `comment`, `comments`, `batch-create`

**Projects** -- `projects` (list), `project-create`, `project-update`

**Cycles** -- `cycles` (list), `cycle-create`

**Initiatives** -- `initiatives` (list), `init-create`

**Labels** -- `labels` (list), `label-create`, `label-delete`

**Workflow States** -- `states` (list), `state-create`

**Documents** -- `docs` (list), `doc-create`

**Users & Auth** -- `users` (list), `me` (current user)

**Webhooks** -- `webhooks` (list), `webhook-create`, `webhook-delete`

**Notifications** -- `notifications` (list), `notif-archive`, `notif-read`

**Meta** -- `teams` (list)

That is 35 commands covering every major Linear entity.

## Installation

```sh
go install github.com/securiter/linear-cli@latest
```

Or download a binary from [Releases](https://github.com/Securiteru/linear-cli/releases).

## Auth

Export your Linear API key:

```sh
export LINEAR_API_KEY=lin_api_...
```

Generate one at **Linear > Settings > API > Personal API keys**.

### With psst (recommended for agents)

psst injects the key at exec time without exposing it to the agent's context:

```sh
psst --global LINEAR_API_KEY -- linear list --team ADI
```

## Usage

### List issues with filters

```sh
linear list --team ADI --status "In Progress" --assignee "Alice" --limit 50
linear list -s "authentication bug" --team ENG
```

### Create an issue

```sh
linear create --title "Fix login redirect" --team ENG --priority 2 --desc "After OAuth..."
```

### Get full issue details

```sh
linear get ENG-142
```

### Update an issue

```sh
linear update ENG-142 --status "Done" --assignee "Bob" --priority 3
linear update ENG-142 --title "New title" --due 2026-06-01T00:00:00Z
linear update ENG-142 --clear-due
```

### Search

```sh
linear search "API rate limit" --limit 10
```

### Comments

```sh
linear comment ENG-142 "Fixed in commit abc123"
linear comments ENG-142
```

### Batch create from stdin

```sh
echo '{"title":"Setup CI","priority":2}
{"title":"Write tests","priority":3}' | linear batch-create --team ENG
```

### JSON output

Most list commands support `--json` for structured output:

```sh
linear list --team ADI --json | jq '.[].title'
linear users --json
```

### Other entities

```sh
linear teams
linear projects --status "Planned"
linear cycles --team ENG
linear initiatives
linear labels --team ENG
linear states --team ENG
linear docs --team ENG --json
linear notifications --limit 50
linear webhooks
linear me
```

### Create/delete entities

```sh
linear project-create --name "Q3 Launch" --team ENG --desc "Ship v2"
linear cycle-create --name "Sprint 42" --team ENG --start 2026-05-01T00:00:00Z --end 2026-05-14T23:59:59Z
linear init-create --name "Platform V3" --target 2026-09-01T00:00:00Z
linear label-create --name "security" --team ENG --color "#ff0000"
linear state-create --name "In Review" --team ENG --type "started" --color "#ffa500"
linear doc-create --title "Architecture RFC" --team ENG --desc "Proposal for..."
linear webhook-create --url https://example.com/webhook --team ENG
```

### Delete/archive

```sh
linear delete ENG-142
linear archive ENG-142
linear unarchive ENG-142
linear webhook-delete <webhook-id>
linear label-delete <label-id>
```

### Notifications

```sh
linear notifications --limit 50
linear notif-archive <notification-id>
linear notif-read <notification-id>
```

## Development

```sh
git clone https://github.com/Securiteru/linear-cli.git
cd linear-cli
go build -o linear .
./linear --help
```

### Project structure

```
main.go              entry point
cmd/
  root.go            root cobra command, auth check
  issues.go          list, search, filter flags
  create.go          issue create + helpers (resolveTeamID, escapeGraphQL)
  get.go             issue get (full detail)
  update.go          issue update (title, status, assignee, priority, labels, due)
  delete.go          delete, archive, unarchive
  comment.go         add/list comments
  batch.go           batch-create from stdin JSON lines
  teams.go           list teams
  labels.go          list/create/delete labels
  states.go          list/create workflow states
  projects.go        list/create/update projects
  cycles.go          list/create cycles
  initiatives.go     list/create initiatives
  documents.go       list/create documents
  users.go           list users
  me.go              current user
  webhooks.go        list/create/delete webhooks
  notifications.go   list/archive/read notifications
api/
  client.go          GraphQL HTTP client
```

## License

MIT
