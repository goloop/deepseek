package deepseek

import (
	"encoding/json"
	"strings"

	"github.com/goloop/ai"
)

// parseError turns a non-success response body into an *ai.APIError. xAI
// reports errors either as {"error":{"message","type","code"}} or as a flat
// {"code":"...","error":"..."} where error is a plain string; both are handled.
func parseError(status int, body []byte) error {
	e := &ai.APIError{
		Status: status,
		Raw:    append(json.RawMessage(nil), body...),
	}

	var obj struct {
		Error struct {
			Message string          `json:"message"`
			Type    string          `json:"type"`
			Code    json.RawMessage `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &obj) == nil && obj.Error.Message != "" {
		e.Message = obj.Error.Message
		e.Type = obj.Error.Type
		e.Code = rawToString(obj.Error.Code)
		return e
	}

	var flat struct {
		Error string          `json:"error"`
		Code  json.RawMessage `json:"code"`
	}
	if json.Unmarshal(body, &flat) == nil {
		e.Message = flat.Error
		e.Code = rawToString(flat.Code)
	}
	return e
}

// rawToString renders a JSON value that may be a string or a number (some
// gateways send a numeric "code") as a plain string.
func rawToString(r json.RawMessage) string {
	s := strings.TrimSpace(string(r))
	if s == "" || s == "null" {
		return ""
	}
	if s[0] == '"' {
		var str string
		if json.Unmarshal(r, &str) == nil {
			return str
		}
	}
	return s
}
