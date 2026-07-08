package deepseek

import (
	"encoding/json"

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
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &obj) == nil && obj.Error.Message != "" {
		e.Message = obj.Error.Message
		e.Type = obj.Error.Type
		e.Code = obj.Error.Code
		return e
	}

	var flat struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if json.Unmarshal(body, &flat) == nil {
		e.Message = flat.Error
		e.Code = flat.Code
	}
	return e
}
