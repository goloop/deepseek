# deepseek - reference

The full reference for the `deepseek` package: the client, the shared
`goloop/ai` model, chat completions (interface and native), streaming, reasoning
content and models.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Creating a client](#creating-a-client)
- [Generate and Stream](#generate-and-stream)
- [Structured output](#structured-output)
- [Native chat completions](#native-chat-completions)
- [Reasoning models](#reasoning-models)
- [Tools and system prompts](#tools-and-system-prompts)
- [Models](#models)
- [Options and errors](#options-and-errors)

## Mental model

`deepseek.Client` implements `ai.Client`, the provider-agnostic contract from
`github.com/goloop/ai`. The shared `Generate` and `Stream` cover the common
ground - chat with tools and streaming - so code written against the interface
runs on any provider.

Provider-specific power lives in native methods: the full `ChatCompletion`
request and model listing. The wire format is chat-completions compatible.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/deepseek"
)
```

## Creating a client

```go
c := deepseek.New(os.Getenv("DEEPSEEK_API_KEY"))

c = deepseek.New(apiKey, deepseek.WithTimeout(30*time.Second))
```

The base URL defaults to `https://api.deepseek.com`. Point `WithBaseURL` at any
compatible endpoint to reuse this client against another gateway.

## Generate and Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    deepseek.ModelChat,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` returns `iter.Seq2[ai.Chunk, error]`: text deltas as chunks with
`Text`, a finished tool call as a chunk with `ToolCall`, and a final chunk with
`Done` and `Usage`.

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		return err
	}
	fmt.Print(chunk.Text)
}
```

## Structured output

`ai.Request.Format` maps onto the provider's own `response_format`, so a request for JSON
is enforced by the provider rather than merely asked for:

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    "the-model",
	Messages: []ai.Message{ai.UserText("Draft SEO fields for this article.")},
	Format: &ai.Format{
		Type:   ai.FormatJSONSchema,
		Name:   "seo",
		Schema: schema,
	},
})

var seo SEO
err = resp.JSON(&seo)
```

`ai.FormatJSON` goes out as `{"type":"json_object"}` and `ai.FormatJSONSchema`
as `{"type":"json_schema", ...}`. Plain JSON mode also appends
`ai.Format.Instruction()` to the system prompt: this wire format rejects
`json_object` unless the word "json" appears in the messages. Your own system
prompt is kept and the instruction follows it; schema mode leaves it untouched.


`ai.Response.Format` is `ai.FormatNative`: this provider enforces every shape
it accepts. Which models support schema mode is the provider's business - there
is no capability table here, so an unsupported pairing is reported by the
provider itself.

## Native chat completions

For provider-only options build a `ChatRequest` and call `ChatCompletion` or
`ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &deepseek.ChatRequest{
	Model:          deepseek.ModelChat,
	Messages:       []deepseek.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

`Tools`, `ToolChoice`, `Temperature`, `TopP`, `MaxTokens`, `Stop`, `N`, `Seed`,
`ResponseFormat` and `User` are all available.

## Reasoning models

`deepseek-reasoner` returns its chain-of-thought in `ReasoningContent`,
separate from the final answer in `Content`:

```go
resp, _ := c.ChatCompletion(ctx, &deepseek.ChatRequest{
	Model:    deepseek.ModelReasoner,
	Messages: []deepseek.ChatMessage{{Role: "user", Content: "9.11 or 9.8?"}},
})
resp.Choices[0].Message.ReasoningContent // the thinking
resp.Choices[0].Message.Content          // the answer
```

`Generate` and `Stream` return only the final answer text. The native
`ChatStream` also streams the reasoning: each `ChatStreamChunk` delta carries
`ReasoningContent` alongside `Content`.

## Tools and system prompts

Tool use and system prompts use the shared `ai` types: `ai.Tool`,
`ai.ToolResult` and a `RoleSystem` message or the `System` field. Tool results
are sent back as `RoleTool` messages whose `ai.ToolResult.ID` matches the
`ai.ToolUse.ID`.

## Models

```go
models, err := c.Models(ctx)
m, err := c.GetModel(ctx, deepseek.ModelChat)
```

## Options and errors

Options: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`.

A non-success response becomes an `*ai.APIError` with `Status`, `Type`, `Code`,
`Message` and the raw body:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// back off
}
```

Requests missing a model or messages fail before the network with
`ai.ErrNoModel` or `ai.ErrNoMessages`.
