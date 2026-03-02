# AGENTS.md

## Project Overview
This repository contains a simple `dotenv` CLI tool written in Go. It loads environment variables from files and command-line arguments before executing a command.

## Project Structure
- `cmd/dotenv/main.go`: The main entry point and implementation of the CLI tool.
- `go.mod`: Dependency definition (requires `github.com/joho/godotenv`).

## Development
- **Language**: Go (v1.15+)
- **Linting**: Run `go vet ./...` to lint the code.
- **Testing**: Run `go test ./...` to run tests (currently no tests are implemented, but the command is standard).

## Dependencies
- `github.com/joho/godotenv`: Used for parsing `.env` files.

## Known Limitations
- Arguments passed with `--key=value` syntax do not support values containing `=`. For example, `--PATH=/usr/bin:/bin` works, but `--CFLAGS=-DDEBUG=1` will fail parsing.
- The tool loads `.env` from the current directory by default, which overrides any command-line arguments with the same key.
