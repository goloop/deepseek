package deepseek

import (
	"fmt"

	"github.com/goloop/ai"
)

// checkHosted refuses a request that asks the provider to run a capability on
// its own side.
//
// This provider's chat endpoint declares function tools only: there is no
// server-side tool to map onto, so there is nothing to run and nothing to cite.
//
// Returning [ai.ErrNoHosted] before the request leaves is the documented
// behavior, not a placeholder: an answer produced without the search that was
// asked for looks exactly like one produced with it, so failing loudly is the
// only way a caller can tell the difference.
func checkHosted(req *ai.Request) error {
	if len(req.Hosted) == 0 {
		return nil
	}
	return fmt.Errorf("%w: this provider's chat endpoint declares function tools only", ai.ErrNoHosted)
}
