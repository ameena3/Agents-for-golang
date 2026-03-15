// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authentication

import (
	"context"
	"fmt"
	"sync"
)

// ConnectionManager manages multiple named OAuth connections.
// The special name "SERVICE_CONNECTION" is treated as the default connection.
type ConnectionManager struct {
	connections map[string]*MsalAuth
	mu          sync.RWMutex
}

// NewConnectionManager creates a new ConnectionManager with a default auth provider
// registered under the name "SERVICE_CONNECTION".
func NewConnectionManager(defaultAuth *MsalAuth) *ConnectionManager {
	cm := &ConnectionManager{
		connections: make(map[string]*MsalAuth),
	}
	if defaultAuth != nil {
		cm.connections["SERVICE_CONNECTION"] = defaultAuth
	}
	return cm
}

// Add registers a named connection. If the name is empty it is registered as
// "SERVICE_CONNECTION" (the default connection).
func (m *ConnectionManager) Add(name string, auth *MsalAuth) {
	if name == "" {
		name = "SERVICE_CONNECTION"
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[name] = auth
}

// GetConnection returns the MsalAuth registered under the given name.
// Returns an error if no connection with that name exists.
func (m *ConnectionManager) GetConnection(name string) (*MsalAuth, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	auth, ok := m.connections[name]
	if !ok {
		return nil, fmt.Errorf("no connection found for %q", name)
	}
	return auth, nil
}

// GetToken retrieves an access token from the named connection for the given resource URL.
func (m *ConnectionManager) GetToken(ctx context.Context, connectionName, resourceURL string) (string, error) {
	auth, err := m.GetConnection(connectionName)
	if err != nil {
		return "", err
	}
	return auth.GetAccessToken(ctx, resourceURL, nil, false)
}

// GetDefaultToken retrieves an access token from the default ("SERVICE_CONNECTION") connection.
func (m *ConnectionManager) GetDefaultToken(ctx context.Context, resourceURL string) (string, error) {
	return m.GetToken(ctx, "SERVICE_CONNECTION", resourceURL)
}
