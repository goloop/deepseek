# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-07-10

### Fixed
- Streamed reasoning is no longer dropped: `ChatCompletionStream` now exposes a
  reasoning model's chain-of-thought in `Delta.ReasoningContent`. The shared
  `Stream` still emits only answer text (reasoning would pollute it).

### Changed
- Require `goloop/ai` v0.2.0 (500 no longer retried; jittered backoff).

## [0.1.0] - 2026-07-09

Initial release, built on the `github.com/goloop/ai` interface.

### Added
- `Client` implementing `ai.Client`: `Generate` and streaming `Stream` over
  chat completions, with tool use and system prompts.
- Native `ChatCompletion` and `ChatCompletionStream` exposing the full chat
  option set; `ChatMessage.ReasoningContent` for reasoning models.
- Models (`Models`, `GetModel`).
- Functional options: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`,
  `WithMaxRetries`, `WithHeader`.
- Retries on 429 and 5xx with backoff; normalized `*ai.APIError` errors.
