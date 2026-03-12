# opsbudget-cli

CLI for [Ping by OpsBudget](https://opsbudget.com) — uptime monitoring from your terminal.

## Install

**macOS (Homebrew):**

```sh
brew install getopsbudget/tap/opsbudget
```

**Linux / macOS (curl):**

```sh
curl -sSL https://get.opsbudget.com | sh
```

## Quick Start

```sh
# Log in (opens browser)
opsbudget login

# Add a monitor
opsbudget ping add https://example.com

# List monitors
opsbudget ping list

# Check status
opsbudget ping status
```

## Commands

| Command | Description |
|---|---|
| `opsbudget login` | Log in via browser |
| `opsbudget ping add <url>` | Add an uptime monitor |
| `opsbudget ping list` | List all monitors |
| `opsbudget ping status` | Show current status |
| `opsbudget ping rm <id\|url>` | Remove a monitor |
| `opsbudget ping history <id\|url>` | Show check history |

### Flags

- `--json` — output as JSON (on any command)
- `--name` — set monitor name (on `ping add`)
- `--method` — HTTP method, default GET (on `ping add`)
- `--interval` — check interval in seconds, 60 or 30 (on `ping add`)
- `--limit` — number of history records (on `ping history`)
- `--force` — skip confirmation (on `ping rm`)

## Configuration

Credentials are stored in `~/.config/opsbudget/credentials.json` after login.

The API base URL can be overridden with the `OPSBUDGET_API_URL` environment variable.

## Docs

Full documentation: [https://opsbudget.com/docs](https://opsbudget.com/docs)

## License

MIT
