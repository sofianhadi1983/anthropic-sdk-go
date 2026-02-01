// Package oauth provides OAuth token authentication for the Anthropic API.
//
// This package enables sending messages using OAuth Bearer tokens, including the
// required beta headers for OAuth authentication.
//
// # Simple Usage
//
// Using environment variable (ANTHROPIC_OAUTH_ACCESS_TOKEN or ANTHROPIC_ACCESS_TOKEN):
//
//	client := anthropic.NewClient(oauth.WithLoadEnv())
//
// Using explicit token:
//
//	client := anthropic.NewClient(oauth.WithAccessToken("your-oauth-token"))
//
// # Full Configuration
//
//	client := anthropic.NewClient(
//	    oauth.WithConfig(oauth.Config{
//	        AccessToken: "your-oauth-token",
//	        Betas: []string{
//	            "oauth-2025-04-20",
//	            "interleaved-thinking-2025-05-14",
//	        },
//	        UserAgent: "my-app/1.0.0",
//	        UseBetaEndpoint: true,
//	    }),
//	)
package oauth

import (
	"net/http"
	"os"
	"strings"

	"github.com/sofianhadi1983/anthropic-sdk-go/internal/requestconfig"
	"github.com/sofianhadi1983/anthropic-sdk-go/option"
)

// DefaultOAuthBetas are the beta features used with OAuth authentication.
var DefaultOAuthBetas = []string{
	"oauth-2025-04-20",
	"interleaved-thinking-2025-05-14",
	"claude-code-20250219",
	"fine-grained-tool-streaming-2025-05-14",
}

// Config holds OAuth authentication configuration.
type Config struct {
	// AccessToken is the OAuth access token for authentication.
	AccessToken string

	// RefreshToken is stored for future token refresh support.
	// Currently not used for automatic refresh.
	RefreshToken string

	// Betas specifies the beta features to enable.
	// Defaults to DefaultOAuthBetas if not set.
	Betas []string

	// UserAgent sets a custom User-Agent header.
	// If empty, the default SDK User-Agent is used.
	UserAgent string

	// UseBetaEndpoint adds ?beta=true query parameter to requests.
	UseBetaEndpoint bool
}

// WithConfig returns a RequestOption for OAuth authentication with full configuration.
//
// Example:
//
//	client := anthropic.NewClient(
//	    oauth.WithConfig(oauth.Config{
//	        AccessToken: "your-oauth-token",
//	        Betas: []string{"oauth-2025-04-20"},
//	        UserAgent: "my-app/1.0.0",
//	        UseBetaEndpoint: true,
//	    }),
//	)
func WithConfig(cfg Config) option.RequestOption {
	if len(cfg.Betas) == 0 {
		cfg.Betas = DefaultOAuthBetas
	}

	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		return rc.Apply(
			option.WithAuthToken(cfg.AccessToken),
			option.WithMiddleware(oauthMiddleware(cfg)),
		)
	})
}

// WithAccessToken returns a RequestOption for OAuth authentication with the given token.
// This uses the default beta features from DefaultOAuthBetas.
//
// Example:
//
//	client := anthropic.NewClient(oauth.WithAccessToken("your-oauth-token"))
func WithAccessToken(accessToken string) option.RequestOption {
	return WithConfig(Config{
		AccessToken: accessToken,
	})
}

// WithLoadEnv returns a RequestOption that loads OAuth configuration from environment variables.
//
// Environment variables:
//   - ANTHROPIC_OAUTH_ACCESS_TOKEN (primary) - the OAuth access token
//   - ANTHROPIC_ACCESS_TOKEN (fallback) - used if ANTHROPIC_OAUTH_ACCESS_TOKEN is not set
//
// Example:
//
//	client := anthropic.NewClient(oauth.WithLoadEnv())
func WithLoadEnv() option.RequestOption {
	accessToken := os.Getenv("ANTHROPIC_OAUTH_ACCESS_TOKEN")
	if accessToken == "" {
		accessToken = os.Getenv("ANTHROPIC_ACCESS_TOKEN")
	}

	return WithConfig(Config{
		AccessToken: accessToken,
	})
}

// oauthMiddleware creates middleware that adds OAuth-specific headers and query parameters.
func oauthMiddleware(cfg Config) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		// Set the anthropic-beta header with OAuth betas
		if len(cfg.Betas) > 0 {
			// Check if there are existing betas to merge with
			existingBetas := r.Header.Get("anthropic-beta")
			if existingBetas != "" {
				// Merge existing betas with OAuth betas, avoiding duplicates
				betaSet := make(map[string]bool)
				for _, b := range strings.Split(existingBetas, ",") {
					betaSet[strings.TrimSpace(b)] = true
				}
				for _, b := range cfg.Betas {
					betaSet[b] = true
				}
				// Rebuild the header
				var allBetas []string
				for b := range betaSet {
					allBetas = append(allBetas, b)
				}
				r.Header.Set("anthropic-beta", strings.Join(allBetas, ","))
			} else {
				r.Header.Set("anthropic-beta", strings.Join(cfg.Betas, ","))
			}
		}

		// Set custom User-Agent if provided
		if cfg.UserAgent != "" {
			r.Header.Set("User-Agent", cfg.UserAgent)
		}

		// Add ?beta=true query parameter if configured
		if cfg.UseBetaEndpoint {
			q := r.URL.Query()
			q.Set("beta", "true")
			r.URL.RawQuery = q.Encode()
		}

		return next(r)
	}
}
