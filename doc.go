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
// # Structured output
//
// ai.Request.Format asks for JSON through the provider's response_format, and
// ai.Response.JSON decodes the reply. The field takes "json_object" and no
// schema type, so ai.FormatJSONSchema also asks for json_object and sends the
// schema itself in the system prompt: valid JSON is enforced, conformance to
// the schema is requested. ai.Response.Format reports which, per format type.
// The instruction is in the prompt either way, since this wire format rejects
// json_object unless the word "json" appears in the messages.
//
// # Hosted capabilities
//
// This provider's chat endpoint declares function tools only, so there is no
// server-side capability to map onto. ai.Hosted is refused with
// ai.ErrNoHosted.
//
// The refusal is the documented behavior, not a gap waiting to be filled
// in silence: an answer produced without the search that was asked for
// looks exactly like one produced with it. A caller who would rather have
// the answer anyway asks again without ai.Request.Hosted.
//
// # Asking what this driver can do
//
// Capabilities describes this driver for the decision taken before a call:
// whether to offer a feature at all, and whether it needs one request or two.
//
//	if ai.SupportsHosted(c, ai.Hosted{Kind: ai.HostedWebSearch}) { ... }
//
// It is a hint and not a permission - support also depends on the model, the
// account and the region - so ai.ErrNoHosted and ai.ErrNoFormat remain the
// source of truth and a caller still handles them. What changes is that a
// refusal the provider only reports as a 400 now arrives as those same
// sentinels, wrapped around the original ai.APIError, so one errors.Is covers
// a limitation this driver knew in advance and one it learned over the wire.
//
// It depends only on goloop/ai and the standard library.
package deepseek
