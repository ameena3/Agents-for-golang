// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package nethttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/microsoft/agents-sdk-go/hosting/core"
)

// ServerConfig configures the HTTP server.
type ServerConfig struct {
	// Port to listen on. Defaults to 3978.
	Port int
	// MessagesPath is the path for the messages endpoint. Defaults to "/api/messages".
	MessagesPath string
	// AllowUnauthenticated disables JWT validation (for local testing).
	AllowUnauthenticated bool
}

// StartAgentProcess starts an HTTP server for the given agent.
// It listens for OS signals and shuts down gracefully.
func StartAgentProcess(ctx context.Context, agent core.Agent, cfg ServerConfig) error {
	if cfg.Port == 0 {
		cfg.Port = 3978
	}
	if cfg.MessagesPath == "" {
		cfg.MessagesPath = "/api/messages"
	}

	adapter := NewCloudAdapter(cfg.AllowUnauthenticated)
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.MessagesPath, MessageHandler(adapter, agent))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-stop:
		case <-ctx.Done():
		}
		_ = srv.Shutdown(context.Background())
	}()

	fmt.Printf("Agent listening on port %d at %s\n", cfg.Port, cfg.MessagesPath)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("nethttp: server error: %w", err)
	}
	return nil
}
