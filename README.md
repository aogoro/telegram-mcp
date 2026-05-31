[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Telegram MCP server (fork)

Fork of [chaindead/telegram-mcp](https://github.com/chaindead/telegram-mcp) with real message sending support.

The server is a bridge between the Telegram API and the AI assistants and is based on the [Model Context Protocol](https://modelcontextprotocol.io).

> [!IMPORTANT]
> Ensure that you have read and understood the [Telegram API Terms of Service](https://core.telegram.org/api/terms) before using this server.
> Any misuse of the Telegram API may result in the suspension of your account.

## Changes from upstream

- **`tg_send`** — sends real messages via `messages.sendMessage` (upstream only saves drafts)
- **`tg_draft`** — saves draft messages (renamed from upstream's `tg_send`)
- **`reply_to`** — optional parameter for `tg_send` to reply to a specific message
- **Proxy fix** — `auth` command now respects `ALL_PROXY` / proxy environment variables
- **`tg_send_photo`** / **`tg_send_file`** — send images (inline preview) and files (documents, exact bytes) to any dialog; both support optional `caption` and `reply_to`

## Capabilities

- [x] Get current account information (`tool: tg_me`)
- [x] List dialogs with optional unread filter (`tool: tg_dialogs`)
- [x] Mark dialog as read (`tool: tg_read`)
- [x] Retrieve messages from specific dialog (`tool: tg_dialog`)
- [x] Send messages to any dialog (`tool: tg_send`) — supports `reply_to` parameter
- [x] Save draft messages to any dialog (`tool: tg_draft`)
- [x] Send a photo/image to any dialog, inline preview (`tool: tg_send_photo`)
- [x] Send a file/document to any dialog, exact bytes (`tool: tg_send_file`)

## Installation

### Homebrew

```bash
# Install
brew install aogoro/tap/telegram-mcp

# Update
brew upgrade aogoro/tap/telegram-mcp
```

### From Releases

Download the latest release for your platform from [Releases](https://github.com/aogoro/telegram-mcp/releases).

#### MacOS

```bash
# For Apple Silicon (M1/M2)
curl -L -o telegram-mcp.tar.gz https://github.com/aogoro/telegram-mcp/releases/latest/download/telegram-mcp_Darwin_arm64.tar.gz

# For Intel Mac (x86_64)
curl -L -o telegram-mcp.tar.gz https://github.com/aogoro/telegram-mcp/releases/latest/download/telegram-mcp_Darwin_x86_64.tar.gz

# Extract and install
sudo tar xzf telegram-mcp.tar.gz -C /usr/local/bin
sudo chmod +x /usr/local/bin/telegram-mcp
rm telegram-mcp.tar.gz
```

#### Linux

```bash
# For x86_64
curl -L -o telegram-mcp.tar.gz https://github.com/aogoro/telegram-mcp/releases/latest/download/telegram-mcp_Linux_x86_64.tar.gz

# For ARM64
curl -L -o telegram-mcp.tar.gz https://github.com/aogoro/telegram-mcp/releases/latest/download/telegram-mcp_Linux_arm64.tar.gz

# Extract and install
sudo tar xzf telegram-mcp.tar.gz -C /usr/local/bin
sudo chmod +x /usr/local/bin/telegram-mcp
rm telegram-mcp.tar.gz
```

### From Source

Requirements: Go 1.24 or later, GOBIN in PATH.

```bash
go install github.com/aogoro/telegram-mcp@latest
```

## Configuration

### Authorization

1. Get the API ID and hash from [Telegram API](https://my.telegram.org/auth)
2. Run:
   ```bash
   telegram-mcp auth --app-id <your-api-id> --api-hash <your-api-hash> --phone <your-phone-number>
   ```
   > If you have 2FA enabled: add `--password <2fa_password>`
   > If you want to override existing session: add `--new`

### Client Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "telegram": {
      "command": "telegram-mcp",
      "env": {
        "TG_APP_ID": "<your-app-id>",
        "TG_API_HASH": "<your-api-hash>"
      }
    }
  }
}
```

### JSON Schema Version

Some MCP clients (e.g. VS Code) do not support JSON Schema Draft 2020-12. Override with `TG_SCHEMA_VERSION`:

| Version | URL |
|---------|-----|
| Draft-07 (recommended for VS Code) | `https://json-schema.org/draft-07/schema#` |
| Draft 2020-12 (default) | `https://json-schema.org/draft/2020-12/schema` |
