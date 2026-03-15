// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authentication

import (
	"testing"
)

func TestNewMsalAuth_MissingTenantAndClient(t *testing.T) {
	_, err := NewMsalAuth(Config{})
	if err == nil {
		t.Fatal("expected error for missing TenantID and ClientID")
	}
}

func TestNewMsalAuth_MissingCredentials(t *testing.T) {
	_, err := NewMsalAuth(Config{TenantID: "tid", ClientID: "cid"})
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

func TestNewMsalAuth_WithSecret(t *testing.T) {
	// This won't actually call AAD, just verifies config creation
	auth, err := NewMsalAuth(Config{
		TenantID:     "test-tenant",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth == nil {
		t.Fatal("expected non-nil auth")
	}
}

func TestNewMsalAuth_WithCustomAuthority(t *testing.T) {
	auth, err := NewMsalAuth(Config{
		TenantID:     "test-tenant",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		Authority:    "https://login.microsoftonline.com/test-tenant",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth == nil {
		t.Fatal("expected non-nil auth")
	}
}

func TestNewConnectionManager(t *testing.T) {
	auth, err := NewMsalAuth(Config{
		TenantID:     "tid",
		ClientID:     "cid",
		ClientSecret: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating auth: %v", err)
	}
	mgr := NewConnectionManager(auth)
	if mgr == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestConnectionManager_Add_And_GetConnection(t *testing.T) {
	auth, err := NewMsalAuth(Config{
		TenantID:     "tid",
		ClientID:     "cid",
		ClientSecret: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating auth: %v", err)
	}

	mgr := NewConnectionManager(nil)
	mgr.Add("MY_CONNECTION", auth)

	got, err := mgr.GetConnection("MY_CONNECTION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil connection")
	}
}

func TestConnectionManager_GetConnection_Missing(t *testing.T) {
	mgr := NewConnectionManager(nil)
	_, err := mgr.GetConnection("NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for nonexistent connection")
	}
}

func TestConnectionManager_DefaultConnection(t *testing.T) {
	auth, err := NewMsalAuth(Config{
		TenantID:     "tid",
		ClientID:     "cid",
		ClientSecret: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating auth: %v", err)
	}

	mgr := NewConnectionManager(auth)
	got, err := mgr.GetConnection("SERVICE_CONNECTION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil default connection")
	}
}
