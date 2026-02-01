package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sofianhadi1983/anthropic-sdk-go"
	"github.com/sofianhadi1983/anthropic-sdk-go/oauth"
	"github.com/sofianhadi1983/anthropic-sdk-go/option"
)

func TestWithAccessToken(t *testing.T) {
	var capturedReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello!"}],"model":"claude-3-5-sonnet-20241022","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":5}}`))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		oauth.WithAccessToken("test-token"),
		option.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(t.Context(), anthropic.MessageNewParams{
		MaxTokens: 256,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello")),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Authorization header
	authHeader := capturedReq.Header.Get("Authorization")
	if authHeader != "Bearer test-token" {
		t.Errorf("expected Authorization header 'Bearer test-token', got '%s'", authHeader)
	}

	// Verify anthropic-beta header contains OAuth beta
	betaHeader := capturedReq.Header.Get("anthropic-beta")
	if !strings.Contains(betaHeader, "oauth-2025-04-20") {
		t.Errorf("expected anthropic-beta header to contain 'oauth-2025-04-20', got '%s'", betaHeader)
	}
}

func TestWithConfig(t *testing.T) {
	var capturedReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello!"}],"model":"claude-3-5-sonnet-20241022","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":5}}`))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		oauth.WithConfig(oauth.Config{
			AccessToken: "custom-token",
			Betas:       []string{"oauth-2025-04-20", "custom-beta"},
			UserAgent:   "test-agent/1.0.0",
		}),
		option.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(t.Context(), anthropic.MessageNewParams{
		MaxTokens: 256,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello")),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify custom User-Agent
	userAgent := capturedReq.Header.Get("User-Agent")
	if userAgent != "test-agent/1.0.0" {
		t.Errorf("expected User-Agent header 'test-agent/1.0.0', got '%s'", userAgent)
	}

	// Verify custom betas
	betaHeader := capturedReq.Header.Get("anthropic-beta")
	if !strings.Contains(betaHeader, "oauth-2025-04-20") || !strings.Contains(betaHeader, "custom-beta") {
		t.Errorf("expected anthropic-beta header to contain custom betas, got '%s'", betaHeader)
	}
}

func TestWithConfigUseBetaEndpoint(t *testing.T) {
	var capturedReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello!"}],"model":"claude-3-5-sonnet-20241022","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":5}}`))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		oauth.WithConfig(oauth.Config{
			AccessToken:     "test-token",
			UseBetaEndpoint: true,
		}),
		option.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(t.Context(), anthropic.MessageNewParams{
		MaxTokens: 256,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello")),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify beta query parameter
	betaParam := capturedReq.URL.Query().Get("beta")
	if betaParam != "true" {
		t.Errorf("expected beta query parameter 'true', got '%s'", betaParam)
	}
}

func TestWithLoadEnv(t *testing.T) {
	// Set environment variable for the test
	os.Setenv("ANTHROPIC_OAUTH_ACCESS_TOKEN", "env-oauth-token")
	defer os.Unsetenv("ANTHROPIC_OAUTH_ACCESS_TOKEN")

	var capturedReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello!"}],"model":"claude-3-5-sonnet-20241022","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":5}}`))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		oauth.WithLoadEnv(),
		option.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(t.Context(), anthropic.MessageNewParams{
		MaxTokens: 256,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello")),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Authorization header uses env var
	authHeader := capturedReq.Header.Get("Authorization")
	if authHeader != "Bearer env-oauth-token" {
		t.Errorf("expected Authorization header 'Bearer env-oauth-token', got '%s'", authHeader)
	}
}

func TestWithLoadEnvFallback(t *testing.T) {
	// Clear primary and set fallback environment variable
	os.Unsetenv("ANTHROPIC_OAUTH_ACCESS_TOKEN")
	os.Setenv("ANTHROPIC_ACCESS_TOKEN", "env-fallback-token")
	defer os.Unsetenv("ANTHROPIC_ACCESS_TOKEN")

	var capturedReq *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello!"}],"model":"claude-3-5-sonnet-20241022","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":10,"output_tokens":5}}`))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		oauth.WithLoadEnv(),
		option.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(t.Context(), anthropic.MessageNewParams{
		MaxTokens: 256,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Hello")),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Authorization header uses fallback env var
	authHeader := capturedReq.Header.Get("Authorization")
	if authHeader != "Bearer env-fallback-token" {
		t.Errorf("expected Authorization header 'Bearer env-fallback-token', got '%s'", authHeader)
	}
}

func TestDefaultOAuthBetas(t *testing.T) {
	// Verify default betas are set correctly
	expected := []string{
		"oauth-2025-04-20",
		"interleaved-thinking-2025-05-14",
		"claude-code-20250219",
		"fine-grained-tool-streaming-2025-05-14",
	}

	if len(oauth.DefaultOAuthBetas) != len(expected) {
		t.Fatalf("expected %d default betas, got %d", len(expected), len(oauth.DefaultOAuthBetas))
	}

	for i, beta := range expected {
		if oauth.DefaultOAuthBetas[i] != beta {
			t.Errorf("expected beta %d to be '%s', got '%s'", i, beta, oauth.DefaultOAuthBetas[i])
		}
	}
}
