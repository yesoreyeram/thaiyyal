//go:build dev
// +build dev

package server

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// In dev mode, we don't embed static files
// Instead, we proxy requests to the Next.js dev server

var nextJSDevURL = "http://localhost:3000"

// getStaticFS returns nil in dev mode as we use a proxy instead
func getStaticFS() (fs.FS, error) {
	// In dev mode, we don't use the embedded filesystem
	return nil, nil
}

// getDevProxy returns a reverse proxy to the Next.js dev server
func getDevProxy() (*httputil.ReverseProxy, error) {
	target, err := url.Parse(nextJSDevURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Customize the proxy director to preserve the original request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	return proxy, nil
}

// isDevMode returns true when built with the dev tag
func isDevMode() bool {
	return true
}
