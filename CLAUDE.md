# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build and install locally
go install ./cmd/codecheck

# Run without installing
go run ./cmd/codecheck <command> [args]

# Build binary
go build ./cmd/codecheck

# Run tests
go test ./...

# Run a specific test
go test ./cmd/codecheck -run TestFunctionName
```

## Architecture

This is a Go CLI tool (`github.com/mentalcaries/codecheck`) with no external dependencies. The entry point is `cmd/codecheck/main.go`.

**Command dispatch** (`commands.go`): A simple hand-rolled dispatcher — `commands.register(name, handlerFn)` maps command names to handler functions. `main.go` registers `review` and `setup`, then routes `os.Args[1]` to the correct handler via `cmds.run(appState, cmd)`.

**Application state** (`main.go`): A `state` struct wraps a `*config.Config` and is passed to every handler. This is the only shared mutable context.

**Config** (`internal/config/config.go`): Persisted to `~/.codecheckconfig.json`. `CheckConfig()` is called at startup — it blocks and prompts the user if no config exists. The stored `DownloadDirectory` is a relative path from `~/`; `Read()` expands it to an absolute path at read time.

**Review flow** (`handler_review.go`): The `handlerReview` function drives the main workflow:
1. Validates and parses the GitHub URL using `ghRegex` (defined in `utils.go`) to extract `user`, `repo`, and optional `branch`
2. Creates `<DownloadDirectory>/<username>/` as the per-user working directory
3. Handles directory conflicts interactively
4. Clones the repo with optional branch support
5. Detects project type: `package.json` present → install deps + `npm run dev`; otherwise → Go `net/http` file server on port 5543 (only if `index.html` exists)
6. Opens VS Code if the `code` CLI is on PATH
7. Blocks on `SIGINT`, then prompts for cleanup

**Regexes** (`utils.go`):
- `ghRegex`: matches both HTTPS and SSH GitHub URLs, optionally capturing a `/tree/<branch>` segment
- `ttProjectRegex`: validates TripleTen student project name format (`<project>-<username>-<commit-hash>`)
