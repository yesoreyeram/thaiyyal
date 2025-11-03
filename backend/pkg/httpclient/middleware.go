package httpclient

import (
	"net/http"
)

// Middleware is a function that wraps an http.RoundTripper
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain creates a chain of middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(base http.RoundTripper) http.RoundTripper {
		// Apply middlewares in reverse order so they execute in the specified order
		for i := len(middlewares) - 1; i >= 0; i-- {
			base = middlewares[i](base)
		}
		return base
	}
}

// authMiddleware adds authentication headers to requests
func authMiddleware(config *Config) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &authRoundTripper{
			next:   next,
			config: config,
		}
	}
}

type authRoundTripper struct {
	next   http.RoundTripper
	config *Config
}

func (t *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add authentication based on type
	switch t.config.Auth.Type {
	case AuthTypeBasic:
		if t.config.Auth.BasicAuth != nil {
			clonedReq.SetBasicAuth(t.config.Auth.BasicAuth.Username, t.config.Auth.BasicAuth.Password.Value())
		}
	case AuthTypeBearer:
		if t.config.Auth.Token != nil {
			clonedReq.Header.Set("Authorization", "Bearer "+t.config.Auth.Token.Token.Value())
		}
	case AuthTypeAPIKey:
		if t.config.Auth.APIKey != nil {
			if t.config.Auth.APIKey.Location == "header" {
				clonedReq.Header.Set(t.config.Auth.APIKey.Key, t.config.Auth.APIKey.Value.Value())
			} else if t.config.Auth.APIKey.Location == "query" {
				q := clonedReq.URL.Query()
				q.Set(t.config.Auth.APIKey.Key, t.config.Auth.APIKey.Value.Value())
				clonedReq.URL.RawQuery = q.Encode()
			}
		}
	}

	return t.next.RoundTrip(clonedReq)
}

// headersMiddleware adds default headers to requests
func headersMiddleware(headers []KeyValue) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &headersRoundTripper{
			next:    next,
			headers: headers,
		}
	}
}

type headersRoundTripper struct {
	next    http.RoundTripper
	headers []KeyValue
}

func (t *headersRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add default headers (don't override existing headers)
	for _, h := range t.headers {
		if clonedReq.Header.Get(h.Key) == "" {
			clonedReq.Header.Add(h.Key, h.Value)
		} else {
			// For headers that already exist, we still add (not set) to support multiple values
			clonedReq.Header.Add(h.Key, h.Value)
		}
	}

	return t.next.RoundTrip(clonedReq)
}

// queryParamsMiddleware adds default query parameters to requests
func queryParamsMiddleware(params []KeyValue) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &queryParamsRoundTripper{
			next:   next,
			params: params,
		}
	}
}

type queryParamsRoundTripper struct {
	next   http.RoundTripper
	params []KeyValue
}

func (t *queryParamsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add default query parameters
	if len(t.params) > 0 {
		q := clonedReq.URL.Query()
		for _, p := range t.params {
			// Add the parameter (supports duplicate keys)
			q.Add(p.Key, p.Value)
		}
		clonedReq.URL.RawQuery = q.Encode()
	}

	return t.next.RoundTrip(clonedReq)
}

// ssrfProtectionMiddleware validates URLs for SSRF protection
func ssrfProtectionMiddleware(config *Config) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &ssrfProtectionRoundTripper{
			next:   next,
			config: config,
		}
	}
}

type ssrfProtectionRoundTripper struct {
	next   http.RoundTripper
	config *Config
}

func (t *ssrfProtectionRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Validate URL before making the request
	if err := validateURL(req.URL.String(), t.config); err != nil {
		return nil, err
	}

	return t.next.RoundTrip(req)
}
