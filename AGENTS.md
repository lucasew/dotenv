# Project Conventions

## General Philosophy
- **Minimize Global State:** Passing the state explicitly as parameters to functions allows them to be purer and easier to test.
- **Centralized Error Reporting:** The project must have a single centralized error-reporting function. All logic flows must pipe unexpected or critical errors through this function. In `cmd/dotenv/main.go`, `handleError` serves this purpose.

## Naming & Directory Structure
- Follow standard Go package layout conventions.
- The main entrypoint is `cmd/dotenv/main.go`.
- Parsing logic for CLI arguments is extracted to `cmd/dotenv/parser.go`.
- Things that change together should live together (colocation).

## Code Style & Rules
- **Formatting:** Always run `go fmt ./...` before committing. Strict CI checks will fail if Go code is not properly formatted.
- **Dependency Management:** Nix (`package.nix`, `default.nix`) is used with `buildGoModule`. When modifying Go package structures or entrypoints, ensure `subPackages` in `package.nix` is correctly configured to instruct Nix where to build, preventing CI build failures.

## Build and Tests
- To build the application: `go build -o build/dotenv ./cmd/dotenv`
- To verify changes locally: `go vet ./... && go test ./...`
- When compiling, target the package directory `./cmd/dotenv` rather than the `main.go` file directly to ensure all package files are included.
- Tooling such as `mise` must be installed via direct standalone binary downloads rather than `curl | sh` scripts. Verification tools run using `mise` (e.g. `mise run ci` or `mise run test`).
