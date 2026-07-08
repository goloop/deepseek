package deepseek

import "github.com/goloop/ai"

// DefaultBaseURL is the DeepSeek API base URL.
const DefaultBaseURL = "https://api.deepseek.com"

// Convenience model identifiers. Any model string is accepted; use Models to
// discover what the account can call.
const (
	ModelChat     = "deepseek-chat"
	ModelReasoner = "deepseek-reasoner"
)

// Client is a DeepSeek API client. It implements [ai.Client] and adds the
// provider's native endpoints. The wire format is chat-completions compatible.
type Client struct {
	opts ai.Options
}

var _ ai.Client = (*Client)(nil)

// New returns a Client for the given API key. Shared options (WithBaseURL,
// WithHTTPClient, WithTimeout, WithMaxRetries, WithHeader) configure it.
func New(apiKey string, opts ...Option) *Client {
	s := settings{}
	for _, o := range opts {
		o(&s)
	}

	o := ai.NewOptions(apiKey, s.aiOpts...)
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}

	return &Client{opts: o}
}
