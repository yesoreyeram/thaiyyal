// Package httpclient provides a standalone HTTP client builder with middleware support.
//
// This package is self-contained and does not depend on any other packages in the project.
// It follows the middleware and roundtripper patterns for extensibility and composability.
//
// # Features
//
//   - Unique immutable UIDs for client identification
//   - Support for duplicate headers and query parameters
//   - Middleware-based authentication (Basic Auth, Bearer Token)
//   - Configurable timeouts and connection pooling
//   - Security-first design with SSRF protection
//   - Thread-safe client registry
//
// # Example Usage
//
//	config := &httpclient.Config{
//	    UID:      "github-api-client",
//	    AuthType: httpclient.AuthTypeBearer,
//	    Token:    "ghp_token",
//	    Timeout:  30 * time.Second,
//	    Headers: []httpclient.KeyValue{
//	        {Key: "Accept", Value: "application/json"},
//	    },
//	}
//
//	client, err := httpclient.New(context.Background(), config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	resp, err := client.Get("https://api.github.com/users/octocat")
package httpclient
