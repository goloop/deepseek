package deepseek

import (
	"context"
	"errors"
	"github.com/goloop/ai"
	"testing"
)

func TestNativeChatRejectsNilRequest(t *testing.T) {
	c := New("")
	if _, err := c.ChatCompletion(context.Background(), nil); !errors.Is(err, ai.ErrNoRequest) {
		t.Fatalf("chat: %v", err)
	}
	n := 0
	for _, err := range c.ChatCompletionStream(context.Background(), nil) {
		n++
		if !errors.Is(err, ai.ErrNoRequest) {
			t.Fatalf("stream: %v", err)
		}
	}
	if n != 1 {
		t.Fatalf("got %d stream errors, want 1", n)
	}
}
