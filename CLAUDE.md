# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the official Go SDK for the Anthropic API. Most code is auto-generated from the OpenAPI spec by Stainless. The `lib/` and `examples/` directories are not modified by the generator.

## Development Commands

```bash
# Build and verify
./scripts/lint              # Build SDK, verify tests compile, build examples

# Run tests (requires mock server)
./scripts/test              # Starts Prism mock server automatically and runs tests
go test ./...               # Run tests directly (requires mock server already running)
go test -run TestName ./... # Run a specific test

# Start mock server manually
./scripts/mock              # Uses OpenAPI spec URL from .stats.yml

# Format code
./scripts/format            # Runs gofmt -s -w
```

## Architecture

### Package Structure

- **Root package (`anthropic`)**: Main client and API types. Import as `github.com/sofianhadi1983/anthropic-sdk-go`
- **`option/`**: Request options for configuring client behavior (headers, retries, timeouts)
- **`bedrock/`**: Amazon Bedrock integration with `bedrock.WithLoadDefaultConfig()` or `bedrock.WithConfig()`
- **`vertex/`**: Google Vertex AI integration with `vertex.WithGoogleAuth()` or `vertex.WithCredentials()`
- **`packages/param/`**: Parameter utilities (`param.Opt[T]`, `param.Null[T]()`, `param.IsOmitted()`)
- **`packages/respjson/`**: Response JSON metadata for field validation
- **`packages/ssestream/`**: Server-sent events streaming support
- **`internal/`**: Internal utilities (JSON encoding, request config, forms, queries)

### Client Services

The `Client` struct exposes these services:
- `client.Messages` - Messages API (main API for Claude interactions)
- `client.Models` - List available models
- `client.Completions` - Legacy completions API
- `client.Beta` - Beta features (files, message batches, skills)

### Key Patterns

**Request Parameters**: Use `omitzero` semantics. Required fields use `json:",required"` tag. Optional primitives use `param.Opt[T]` with constructors like `anthropic.String()`, `anthropic.Int()`.

**Request Unions**: Struct with `Of*` prefixed fields (e.g., `OfTool`, `OfCat`). Only one can be non-zero.

**Response Unions**: Flattened struct with all variant fields. Use `.AsAny()` for type switching or `.AsFooVariant()` methods.

**Response Validation**: Use `.JSON.FieldName.Valid()` to check if optional fields were present/non-null.

**Streaming**: Use `client.Messages.NewStreaming()` returning an iterator. Call `message.Accumulate(event)` to build complete response.

## Running Examples

```bash
go run ./examples/<example-name>
```

Available examples: message, message-streaming, tools, tools-streaming, multimodal, bedrock, vertex, structured-outputs, etc.
