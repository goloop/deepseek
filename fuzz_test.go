package deepseek

import (
	"errors"
	"testing"

	"github.com/goloop/ai"
)

// FuzzParseError checks that error parsing never panics on arbitrary response
// bodies and always yields an *ai.APIError that preserves the status and the
// raw body, across both the nested and flat error shapes.
func FuzzParseError(f *testing.F) {
	for _, s := range []string{
		`{"error":{"message":"bad","type":"invalid_request_error","code":"x"}}`,
		`{"error":{"message":"bad","code":123}}`,
		`{"error":"flat string","code":"e"}`,
		`{"message":"flat message","type":"t","code":42}`,
		`not json`,
		``,
		`{"error":{"code":null}}`,
		`{"code":"  spaced  "}`,
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, body string) {
		err := parseError(503, []byte(body))
		var apiErr *ai.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("not an *ai.APIError: %T", err)
		}
		if apiErr.Status != 503 {
			t.Fatalf("status = %d", apiErr.Status)
		}
		if string(apiErr.Raw) != body {
			t.Fatalf("raw body not preserved: %q != %q", apiErr.Raw, body)
		}
		_ = rawToString([]byte(body))
	})
}
