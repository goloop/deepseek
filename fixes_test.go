package deepseek

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/goloop/ai"
)

// BUG-01: tool calls flushed even without a trailing finish_reason.
func TestFixStreamToolCallWithoutFinishReason(t *testing.T) {
	events := []string{
		`data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1",` +
			`"function":{"name":"lookup","arguments":"{\"q\":42}"}}]}}]}`, ``,
		`data: [DONE]`, ``,
	}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for _, line := range events {
			io.WriteString(w, line+"\n")
		}
	})
	defer done()

	var call *ai.ToolUse
	for chunk, err := range c.Stream(context.Background(), &ai.Request{
		Model: "m", Messages: []ai.Message{ai.UserText("hi")},
	}) {
		if err != nil {
			t.Fatal(err)
		}
		if chunk.ToolCall != nil {
			call = chunk.ToolCall
		}
	}
	if call == nil || call.Name != "lookup" || string(call.Input) != `{"q":42}` {
		t.Fatalf("tool call lost: %+v", call)
	}
}

// BUG-02: native ChatCompletion must not mutate the caller's request.
func TestFixChatRequestNotMutated(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"model":"m","choices":[{"index":0,`+
			`"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`)
	})
	defer done()

	req := &ChatRequest{Model: "m", Messages: []ChatMessage{{Role: "user", Content: "hi"}}}
	if _, err := c.ChatCompletion(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if req.Stream || req.StreamOptions != nil {
		t.Error("ChatCompletion mutated the caller's request")
	}
}

// BUG-03: content returned as an array of parts must not be dropped.
func TestFixContentPartsResponse(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"model":"m","choices":[{"index":0,"message":{"role":"assistant",`+
			`"content":[{"type":"text","text":"Hello "},{"type":"text","text":"world"}]},`+
			`"finish_reason":"stop"}]}`)
	})
	defer done()

	resp, err := c.Generate(context.Background(), &ai.Request{
		Model: "m", Messages: []ai.Message{ai.UserText("hi")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() != "Hello world" {
		t.Errorf("text = %q", resp.Text())
	}
}

// BUG-04: a numeric "code" must not break error parsing.
func TestFixParseErrorNumericCode(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error":{"message":"bad","type":"invalid_request_error","code":429}}`)
	})
	defer done()

	_, err := c.Generate(context.Background(), &ai.Request{
		Model: "m", Messages: []ai.Message{ai.UserText("hi")},
	})
	var apiErr *ai.APIError
	if !errors.As(err, &apiErr) || apiErr.Message != "bad" {
		t.Fatalf("message lost: %v", err)
	}
	if apiErr.Code != "429" {
		t.Errorf("code = %q, want 429", apiErr.Code)
	}
}
