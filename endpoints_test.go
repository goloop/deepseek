package deepseek

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestChatCompletionNative(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"model":"m","choices":[{"index":0,`+
			`"message":{"role":"assistant","content":"hi","reasoning_content":"thinking"},`+
			`"finish_reason":"stop"}]}`)
	})
	defer done()

	resp, err := c.ChatCompletion(context.Background(), &ChatRequest{
		Model:    ModelReasoner,
		Messages: []ChatMessage{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	msg := resp.Choices[0].Message
	if msg.Content != "hi" || msg.ReasoningContent != "thinking" {
		t.Errorf("message = %+v", msg)
	}
}

func TestChatCompletionStreamNative(t *testing.T) {
	events := []string{
		`data: {"choices":[{"index":0,"delta":{"content":"a"}}]}`, ``,
		`data: {"choices":[{"index":0,"delta":{"content":"b"}}]}`, ``,
		`data: [DONE]`, ``,
	}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for _, line := range events {
			io.WriteString(w, line+"\n")
		}
	})
	defer done()

	var text strings.Builder
	for chunk, err := range c.ChatCompletionStream(context.Background(), &ChatRequest{
		Model: "m", Messages: []ChatMessage{{Role: "user", Content: "hi"}},
	}) {
		if err != nil {
			t.Fatal(err)
		}
		for _, ch := range chunk.Choices {
			text.WriteString(ch.Delta.Content)
		}
	}
	if text.String() != "ab" {
		t.Errorf("text = %q", text.String())
	}
}

func TestModels(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/models/deepseek-chat"):
			io.WriteString(w, `{"id":"deepseek-chat","object":"model"}`)
		default:
			io.WriteString(w, `{"data":[{"id":"deepseek-chat","object":"model"}]}`)
		}
	})
	defer done()

	ctx := context.Background()
	if models, err := c.Models(ctx); err != nil || len(models) != 1 {
		t.Fatalf("models: %v %+v", err, models)
	}
	m, err := c.GetModel(ctx, "deepseek-chat")
	if err != nil || m.ID != "deepseek-chat" {
		t.Fatalf("get model: %v %+v", err, m)
	}
}
