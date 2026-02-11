# postcode-pt-cli

[![CI](https://github.com/RobertoCCC/postcode-pt-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/RobertoCCC/postcode-pt-cli/actions/workflows/ci.yml)
[![Release](https://github.com/RobertoCCC/postcode-pt-cli/actions/workflows/release.yml/badge.svg)](https://github.com/RobertoCCC/postcode-pt-cli/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/RobertoCCC/postcode-pt-cli.svg)](https://pkg.go.dev/github.com/RobertoCCC/postcode-pt-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/RobertoCCC/postcode-pt-cli)](https://goreportcard.com/report/github.com/RobertoCCC/postcode-pt-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`pcpt` is a small command-line client for the [postcode-pt](https://github.com/RobertoCCC/postcode-pt) API — look up Portuguese postal codes and browse the district / municipality hierarchy without leaving the terminal.

```text
$ pcpt 1100-038
1100-038  LISBOA
  street       Rua do Arsenal
  locality     Lisboa (21696)
  municipality Lisboa (1106)
  district     Lisboa (11)
```

## Install

### Pre-built binaries

Grab the archive for your platform from the [latest release](https://github.com/RobertoCCC/postcode-pt-cli/releases/latest) and drop `pcpt` somewhere on your `PATH`.

### `go install`

```bash
go install github.com/RobertoCCC/postcode-pt-cli@latest
```

The binary lands in `$(go env GOBIN)` (or `$(go env GOPATH)/bin`). Rename it to `pcpt` if you prefer the short name:

```bash
mv "$(go env GOPATH)/bin/postcode-pt-cli" "$(go env GOPATH)/bin/pcpt"
```

### Build from source

```bash
git clone https://github.com/RobertoCCC/postcode-pt-cli.git
cd postcode-pt-cli
go build -o pcpt .
```

## Usage

### Look up a postal code

```bash
pcpt 1100-038      # canonical form
pcpt 1100038       # 7 digits also work
```

### List districts

```bash
pcpt districts
```

### List municipalities of a district

```bash
pcpt district 11   # Lisboa
```

### JSON output (pipe to `jq`)

```bash
pcpt --json 1100-038 | jq '.[0].district'
```

### Point at a different API instance

```bash
pcpt --api-url http://localhost:8000/v1 1100-038
```

## Flags

| Flag         | Description                                              |
| ------------ | -------------------------------------------------------- |
| `--api-url`  | Override the API base URL                                |
| `--json`     | Emit raw JSON instead of the human-readable layout       |
| `--no-color` | Disable ANSI colors (also honours `NO_COLOR`)            |

Colors are auto-disabled when stdout is not a TTY, so piping to other tools is safe.

## Project layout

```
.
├── main.go
├── cmd/                 # cobra commands
├── internal/
│   ├── api/             # HTTP client + types
│   └── render/          # TTY-aware output formatter
├── .goreleaser.yaml     # cross-platform release config
└── .github/workflows/   # CI + release pipelines
```

## Development

```bash
go test ./...        # unit + httptest
go vet ./...
go build ./...
```

CI runs the same matrix on Linux, macOS and Windows.

## Releases

Tagged pushes (`vX.Y.Z`) trigger [GoReleaser](https://goreleaser.com), which cross-compiles binaries for Linux / macOS / Windows on amd64 and arm64, publishes them to GitHub Releases, and ships a checksum file.

## Related projects

- [postcode-pt](https://github.com/RobertoCCC/postcode-pt) — the FastAPI service this CLI talks to
- [postcode-pt-web](https://github.com/RobertoCCC/postcode-pt-web) — Next.js frontend backed by the same API

## License

[MIT](LICENSE)
