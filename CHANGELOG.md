# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-08-11

First stable release, on `ai` v1.0.0.

### Added
- `ai.Request.Hosted` is answered with `ai.ErrNoHosted` before the request
  leaves, because the provider's chat endpoint declares function tools only, so there is no
  server-side capability to map onto.
  That refusal is the documented behavior rather than a gap left in silence:
  an answer produced without the search that was asked for looks exactly like
  one produced with it, so failing loudly is the only way a caller can tell
  them apart. A caller who would rather have the answer anyway asks again
  without `Hosted`.
- Tests pin it, including that nothing is sent to the provider.

## [0.2.3] - 2026-08-05

### Fixed
- A request carrying `ai.FormatJSONSchema` no longer fails outright. The driver
  was sending `response_format: {"type":"json_schema", ...}`, copied from the
  rest of this wire format's family, but this provider's field accepts only
  `text` and `json_object` - so every schema request was rejected by the
  endpoint while `ai.Response.Format` claimed `ai.FormatNative`. A schema now
  asks for `json_object`, which does constrain the reply to valid JSON, and
  travels to the model in the system prompt where `ai.Format.Instruction`
  already spells it out.
- `ai.Response.Format` tells the two apart: `ai.FormatNative` for
  `ai.FormatJSON`, since the provider enforces valid JSON, and
  `ai.FormatEmulated` for `ai.FormatJSONSchema`, since the schema is a request
  to the model rather than a constraint on it. Code that needs schema
  conformance guaranteed can now see that it is not.

## [0.2.2] - 2026-08-05

### Documentation
- The package documentation describes how `ai.Request.Format` reaches this
  provider, so it is on the first page a reader sees rather than only in the
  reference.

## [0.2.1] - 2026-08-05

Released by mistake: it carries this changelog entry and nothing else. The
documentation it describes landed in 0.2.2.

## [0.2.0] - 2026-08-05

### Added
- `ai.Request.Format` is mapped onto the provider's own `response_format`, so a request
  for JSON is enforced by the provider rather than merely asked for, and
  `ai.Response.JSON` decodes the reply. `Response.Format` reports
  `ai.FormatNative`. Until now only this package's native request type could
  ask for JSON, so callers going through the provider-agnostic interface had
  to strip code fences from the reply by hand.
- Plain JSON mode also appends `ai.Format.Instruction()` to the system prompt,
  because this wire format rejects `json_object` unless the word "json" appears
  in the messages. The caller's own system prompt is kept and the instruction
  follows it; schema mode leaves the prompt untouched.

### Changed
- Requires `github.com/goloop/ai` v0.4.0.

## [0.1.2] - 2026-07-10

### Documentation
- `DOC.md`/`DOC.UK.md` note that the native `ChatStream` streams
  `ReasoningContent` in each delta alongside the content.

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
