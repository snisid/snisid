package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/snisid/platform/internal/platform/logger"
)

// NewReverseProxy creates a customized httputil.ReverseProxy.
func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the transport for timeouts and keep-alives
	proxy.Transport = &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = targetURL.Host

		// Inject Correlation ID if missing
		if req.Header.Get("X-Correlation-ID") == "" {
			if cid, ok := req.Context().Value("correlation_id").(string); ok && cid != "" {
				req.Header.Set("X-Correlation-ID", cid)
			}
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error(context.Background(), "proxy error", err)
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error": "bad gateway"}`))
	}

	return proxy, nil
}

// ProxyHandler returns an http.Handler that routes to the given target.
func ProxyHandler(target string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy, err := NewReverseProxy(target)
		if err != nil {
			http.Error(w, "invalid target url", http.StatusInternalServerError)
			return
		}

		// Add timeout context to proxy
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		proxy.ServeHTTP(w, r.WithContext(ctx))
	})
}
