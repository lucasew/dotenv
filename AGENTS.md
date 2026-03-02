# Agent Guidelines

## General Context
This repository provides a Go CLI tool (`dotenv`) for launching commands with environment variables loaded from specific files or arguments.
The project relies on standard Go tools (`go test`, `go vet`) and uses `mise` for tool management and task execution.

## Directory Pointers
- `cmd/dotenv/` -> Entrypoint for the CLI application.
- `cmd/dotenv/main.go` -> Main logic and execution flow.
- `cmd/dotenv/parser.go` -> Logic for parsing arguments and files (to be extracted here).

## Coding Guidelines
- **Go Conventions**: Adhere strictly to Go standard project layout and idiomatic Go formatting.
- **Global State**: Minimize the use of global state. Inject variables explicitly where possible to maintain modularity.
- **Error Handling**: Use the centralized error handling function `handleError` in `main.go`. All unrecoverable errors should funnel through this single point rather than using inline `os.Exit` or `log.Fatal` elsewhere.
- **Single Responsibility**: Each file should have a distinct, well-defined responsibility. Code extraction and separation of concerns are heavily encouraged based on Robert C. Martin's Single Responsibility Principle.

## Testing
- Tests should live alongside the code they test (e.g. `cmd/dotenv/parser_test.go`).
- Always run automated checks (`mise run ci` or `go test ./...` and `go vet ./...`) before confirming changes.
