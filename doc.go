// Package deepseek is a client for the DeepSeek API, built on the goloop/ai
// interface.
//
// The Client implements ai.Client, so Generate and Stream work the same as
// with any other goloop AI provider. On top of that it exposes the native
// chat completions endpoint with its full options and model listing. The wire
// format is chat-completions compatible; reasoning models return their
// chain-of-thought in ChatMessage.ReasoningContent (and, when streaming,
// ChatStreamChunk's Delta.ReasoningContent).
//
//	c := deepseek.New(os.Getenv("DEEPSEEK_API_KEY"))
//	resp, err := c.Generate(ctx, &ai.Request{
//	    Model:    deepseek.ModelChat,
//	    Messages: []ai.Message{ai.UserText("Say hello in one word.")},
//	})
//
// It depends only on goloop/ai and the standard library.
package deepseek
