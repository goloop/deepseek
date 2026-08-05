package deepseek

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/goloop/ai"
)

// askFor builds a minimal request carrying f.
func askFor(f *ai.Format) *ai.Request {
	return &ai.Request{
		Model:    "the-model",
		System:   "You are terse.",
		Messages: []ai.Message{ai.UserText("Describe this article.")},
		Format:   f,
	}
}

// TestResponseFormatJSON checks plain JSON mode goes out as the provider's own
// json_object, with the word this wire format insists on in the prompt.
func TestResponseFormatJSON(t *testing.T) {
	cr, err := (&Client{}).chatRequest(askFor(&ai.Format{Type: ai.FormatJSON}), false)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(cr.ResponseFormat); got != `{"type":"json_object"}` {
		t.Errorf("response_format = %s", got)
	}

	system, _ := cr.Messages[0].Content.(string)
	if !strings.HasPrefix(system, "You are terse.") {
		t.Errorf("caller's system prompt was lost: %q", system)
	}
	if !strings.Contains(strings.ToLower(system), "json") {
		t.Errorf("system prompt does not mention json: %q", system)
	}
}

// TestResponseFormatJSONSchema pins what a schema turns into here. The field
// takes "text" or "json_object" and nothing else, so a request carrying
// {"type":"json_schema"} would be rejected by the endpoint every time - the
// schema asks for json_object and reaches the model through the prompt.
func TestResponseFormatJSONSchema(t *testing.T) {
	schema := json.RawMessage(`{"type":"object","properties":{"a":{"type":"string"}}}`)
	cr, err := (&Client{}).chatRequest(askFor(&ai.Format{
		Type:   ai.FormatJSONSchema,
		Name:   "seo",
		Schema: schema,
		Strict: true,
	}), false)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(cr.ResponseFormat); got != `{"type":"json_object"}` {
		t.Errorf("response_format = %s, want json_object", got)
	}
	if strings.Contains(string(cr.ResponseFormat), "json_schema") {
		t.Errorf("response_format carries a type this endpoint rejects: %s",
			cr.ResponseFormat)
	}

	system, _ := cr.Messages[0].Content.(string)
	if !strings.HasPrefix(system, "You are terse.") {
		t.Errorf("caller's system prompt was lost: %q", system)
	}
	if !strings.Contains(system, string(schema)) {
		t.Errorf("the schema never reached the model: %q", system)
	}
}

// TestResponseFormatDefaults checks a request that asks for nothing still
// looks exactly as it did before the field existed.
func TestResponseFormatDefaults(t *testing.T) {
	for _, f := range []*ai.Format{nil, {Type: ai.FormatText}} {
		cr, err := (&Client{}).chatRequest(askFor(f), false)
		if err != nil {
			t.Fatal(err)
		}
		if cr.ResponseFormat != nil {
			t.Errorf("response_format = %s, want it absent", cr.ResponseFormat)
		}
		if system, _ := cr.Messages[0].Content.(string); system != "You are terse." {
			t.Errorf("system prompt was changed: %q", system)
		}
		if got := formatMode(f); got != ai.FormatNone {
			t.Errorf("formatMode = %s, want none", got)
		}
	}
}

// TestResponseFormatRejectsUnknown checks a format this driver cannot render
// fails here, not as a puzzling error from the provider.
func TestResponseFormatRejectsUnknown(t *testing.T) {
	_, err := (&Client{}).chatRequest(askFor(&ai.Format{Type: ai.FormatType(99)}), false)
	if !errors.Is(err, ai.ErrBadFormat) {
		t.Errorf("err = %v, want %v", err, ai.ErrBadFormat)
	}

	_, err = (&Client{}).chatRequest(askFor(&ai.Format{Type: ai.FormatJSONSchema}), false)
	if !errors.Is(err, ai.ErrNoSchema) {
		t.Errorf("err = %v, want %v", err, ai.ErrNoSchema)
	}
}

// TestFormatMode records what this driver promises, and where it stops. Valid
// JSON the provider enforces; conformance to a schema it does not, because the
// schema is in the prompt. A caller that needs the schema guaranteed has to be
// able to see that.
func TestFormatMode(t *testing.T) {
	cases := map[ai.FormatMode][]*ai.Format{
		ai.FormatNative: {
			{Type: ai.FormatJSON},
		},
		ai.FormatEmulated: {
			{Type: ai.FormatJSONSchema, Schema: json.RawMessage(`{"type":"object"}`)},
		},
		ai.FormatNone: {
			nil,
			{Type: ai.FormatText},
			{Type: ai.FormatType(99)},
		},
	}
	for want, formats := range cases {
		for _, f := range formats {
			if got := formatMode(f); got != want {
				t.Errorf("formatMode(%v) = %s, want %s", f, got, want)
			}
		}
	}
}

// TestStreamAsksTheSame checks streaming carries the format exactly as the
// single-shot path does.
func TestStreamAsksTheSame(t *testing.T) {
	f := &ai.Format{Type: ai.FormatJSON}
	direct, err := (&Client{}).chatRequest(askFor(f), false)
	if err != nil {
		t.Fatal(err)
	}
	streamed, err := (&Client{}).chatRequest(askFor(f), true)
	if err != nil {
		t.Fatal(err)
	}
	if string(direct.ResponseFormat) != string(streamed.ResponseFormat) ||
		len(direct.Messages) != len(streamed.Messages) {
		t.Errorf("stream and direct requests differ")
	}
}
