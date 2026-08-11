[![deps.dev](https://img.shields.io/badge/deps.dev-insights-4c8dbc)](https://deps.dev/go/github.com%2Fgoloop%2Fdeepseek) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/deepseek/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/deepseek) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# deepseek

`deepseek` is a Go client for the DeepSeek API. It implements the
`github.com/goloop/ai` interface, so it looks and works like every other goloop
AI provider, and exposes the native chat-completions endpoint with its full
options on top.

## Features
- Chat completions: `Generate` for a single response, `Stream` for
  token-by-token output through `iter.Seq2`.
- Tool use (function calling) and system prompts.
- Native `ChatCompletion` and `ChatCompletionStream` with the full option set;
  reasoning models return their chain-of-thought in `ReasoningContent`.
- Model listing.
- Retries on 429 and 5xx with backoff; normalized, typed API errors.
- Depends only on `github.com/goloop/ai` and the standard library.
- Structured output: `ai.Format` asks for JSON through the provider's
  `response_format`; a schema travels in the prompt, since the field has no
  schema type. Read the reply with `resp.JSON(&v)`.
- Hosted capabilities: `ai.Request.Hosted` is refused with `ai.ErrNoHosted`,
  because this chat endpoint declares function tools only.

## Installation

```sh
go get github.com/goloop/deepseek
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goloop/ai"
	"github.com/goloop/deepseek"
)

func main() {
	c := deepseek.New(os.Getenv("DEEPSEEK_API_KEY"))

	resp, err := c.Generate(context.Background(), &ai.Request{
		Model:    deepseek.ModelChat,
		Messages: []ai.Message{ai.UserText("Say hello in one word.")},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}
```

## Streaming

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		break
	}
	fmt.Print(chunk.Text)
	if chunk.Done && chunk.Usage != nil {
		fmt.Printf("\n[%d in / %d out]\n",
			chunk.Usage.InputTokens, chunk.Usage.OutputTokens)
	}
}
```

## Reasoning models

`deepseek-reasoner` returns its chain-of-thought alongside the answer. Read it
from a native response:

```go
resp, _ := c.ChatCompletion(ctx, &deepseek.ChatRequest{
	Model:    deepseek.ModelReasoner,
	Messages: []deepseek.ChatMessage{{Role: "user", Content: "9.11 or 9.8, which is bigger?"}},
})
fmt.Println(resp.Choices[0].Message.ReasoningContent) // the thinking
fmt.Println(resp.Choices[0].Message.Content)          // the answer
```

## Documentation

Full reference: **[DOC.md](DOC.md)** (Ukrainian: **[DOC.UK.md](DOC.UK.md)**).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT - see [LICENSE](LICENSE).
