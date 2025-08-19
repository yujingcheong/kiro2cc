# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI tool called `kiro2cc` that manages Kiro authentication tokens and provides an Anthropic API proxy service. The tool acts as a bridge between Anthropic API requests and AWS CodeWhisperer, translating requests and responses between the two formats.

## Quick Start Guide

### Prerequisites
1. Install and login to Kiro desktop application from [kiro.dev](https://kiro.dev)
2. Use Kiro AI features at least once to generate authentication tokens
3. Ensure Go 1.23+ is installed (if building from source)

### Getting Started
1. **Download**: Get the latest release from GitHub releases page
2. **Verify tokens**: `./kiro2cc read` to check if Kiro tokens exist
3. **Refresh if needed**: `./kiro2cc refresh` to update expired tokens  
4. **Start server**: `./kiro2cc server` to run the proxy on port 8080
5. **Test**: Send requests to `http://localhost:8080/v1/messages`

For detailed instructions, see the main README.md file.

## Build and Development Commands

```bash
# Build the application
go build -o kiro2cc main.go

# Run tests
go test ./...

# Run specific test in parser package
go test ./parser -v

# Run the application
./kiro2cc [command]
```

## Application Commands

- `./kiro2cc read` - Read and display token information
- `./kiro2cc refresh` - Refresh the access token using refresh token
- `./kiro2cc export` - Export environment variables for other tools
- `./kiro2cc server [port]` - Start HTTP proxy server (default port 8080)

## Architecture

### Core Components

1. **Token Management** (`main.go:337-486`)
   - Reads tokens from `~/.aws/sso/cache/kiro-auth-token.json`
   - Handles token refresh via Kiro auth service
   - Cross-platform environment variable export

2. **API Translation** (`main.go:232-303`)
   - Converts Anthropic API requests to CodeWhisperer format
   - Maps model names via `ModelMap` (line 218-221)
   - Handles conversation history and system messages

3. **HTTP Proxy Server** (`main.go:514-893`)
   - Serves on `/v1/messages` endpoint
   - Supports both streaming and non-streaming requests
   - Automatic token refresh on 403 errors

4. **Response Parser** (`parser/sse_parser.go`)
   - Parses binary CodeWhisperer responses
   - Converts to Anthropic-compatible SSE events
   - Handles tool use and text content blocks

### Key Data Structures

- `AnthropicRequest` - Incoming API requests
- `CodeWhispererRequest` - Outgoing AWS requests  
- `TokenData` - Authentication token storage
- `SSEEvent` - Streaming response events

### Request Flow

1. Client sends Anthropic API request to `/v1/messages`
2. Server reads token from filesystem
3. Request converted to CodeWhisperer format
4. Proxied to AWS CodeWhisperer API
5. Response parsed and converted back to Anthropic format
6. Streamed or returned as JSON to client

## Development Notes

- Uses hardcoded proxy `127.0.0.1:9000` for AWS requests
- Model mapping required between Anthropic and CodeWhisperer model IDs
- Response files saved as `msg_[timestamp]response.raw` for debugging
- Automatic token refresh on authentication failures