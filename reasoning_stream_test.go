package deepseek

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestChatCompletionStreamReasoning verifies the streamed chain-of-thought of a
// reasoning model is surfaced in Delta.ReasoningContent (it was previously
// dropped because the stream delta had no reasoning_content field).
func TestChatCompletionStreamReasoning(t *testing.T) {
	events := []string{
		`data: {"choices":[{"index":0,"delta":{"reasoning_content":"let me"}}]}`,
		``,
		`data: {"choices":[{"index":0,"delta":{"reasoning_content":" think"}}]}`,
		``,
		`data: {"choices":[{"index":0,"delta":{"content":"answer"}}]}`,
		``,
		`data: [DONE]`,
		``,
	}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for _, line := range events {
			io.WriteString(w, line+"\n")
		}
	})
	defer done()

	var reasoning, content strings.Builder
	for chunk, err := range c.ChatCompletionStream(context.Background(), &ChatRequest{
		Model:    "deepseek-reasoner",
		Messages: []ChatMessage{{Role: "user", Content: "hi"}},
	}) {
		if err != nil {
			t.Fatal(err)
		}
		for _, ch := range chunk.Choices {
			reasoning.WriteString(ch.Delta.ReasoningContent)
			content.WriteString(ch.Delta.Content)
		}
	}
	if reasoning.String() != "let me think" {
		t.Errorf("reasoning = %q, want %q", reasoning.String(), "let me think")
	}
	if content.String() != "answer" {
		t.Errorf("content = %q", content.String())
	}
}
