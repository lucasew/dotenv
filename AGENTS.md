# Project Conventions

## Where To Find Things
- `cmd/dotenv/main.go` -> Entrypoint and centralized error handling (`handleError`).
- `cmd/dotenv/parser.go` -> Parsing logic (if extracted).

## Development
- Use `go fmt ./...` and `go vet ./...` before submitting changes.
- Errors must never be silently swallowed.
- Use explicit error checks for expected errors (like `os.IsNotExist` for optional files).
- Always use explicit file staging (`git add <file>`) instead of `git add -A`.
