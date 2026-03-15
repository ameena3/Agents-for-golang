// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	// Clear all relevant env vars to get default zero values.
	vars := []string{
		"MICROSOFT_APP_ID", "AZURE_CLIENT_ID",
		"MICROSOFT_APP_PASSWORD", "AZURE_CLIENT_SECRET",
		"MICROSOFT_APP_TENANT_ID", "AZURE_TENANT_ID",
		"SERVICE_URL",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}

	cfg := LoadFromEnv()
	if cfg == nil {
		t.Fatal("LoadFromEnv returned nil")
	}
	if cfg.AppID != "" {
		t.Errorf("AppID: got %q, want empty string", cfg.AppID)
	}
	if cfg.ClientID != "" {
		t.Errorf("ClientID: got %q, want empty string", cfg.ClientID)
	}
	if cfg.AppPassword != "" {
		t.Errorf("AppPassword: got %q, want empty string", cfg.AppPassword)
	}
	if cfg.ClientSecret != "" {
		t.Errorf("ClientSecret: got %q, want empty string", cfg.ClientSecret)
	}
	if cfg.TenantID != "" {
		t.Errorf("TenantID: got %q, want empty string", cfg.TenantID)
	}
	if cfg.ServiceURL != "" {
		t.Errorf("ServiceURL: got %q, want empty string", cfg.ServiceURL)
	}
	if cfg.Connections == nil {
		t.Error("Connections map should be non-nil")
	}
}

func TestLoadFromEnv_MicrosoftAppID(t *testing.T) {
	os.Unsetenv("AZURE_CLIENT_ID")
	os.Setenv("MICROSOFT_APP_ID", "my-app-id")
	defer os.Unsetenv("MICROSOFT_APP_ID")

	cfg := LoadFromEnv()
	if cfg.AppID != "my-app-id" {
		t.Errorf("AppID: got %q, want %q", cfg.AppID, "my-app-id")
	}
	if cfg.ClientID != "my-app-id" {
		t.Errorf("ClientID: got %q, want %q", cfg.ClientID, "my-app-id")
	}
}

func TestLoadFromEnv_AzureClientIDFallback(t *testing.T) {
	os.Unsetenv("MICROSOFT_APP_ID")
	os.Setenv("AZURE_CLIENT_ID", "azure-client-id")
	defer os.Unsetenv("AZURE_CLIENT_ID")

	cfg := LoadFromEnv()
	if cfg.AppID != "azure-client-id" {
		t.Errorf("AppID fallback: got %q, want %q", cfg.AppID, "azure-client-id")
	}
}

func TestLoadFromEnv_PrimaryTakesPrecedenceOverFallback(t *testing.T) {
	os.Setenv("MICROSOFT_APP_ID", "primary")
	os.Setenv("AZURE_CLIENT_ID", "fallback")
	defer func() {
		os.Unsetenv("MICROSOFT_APP_ID")
		os.Unsetenv("AZURE_CLIENT_ID")
	}()

	cfg := LoadFromEnv()
	if cfg.AppID != "primary" {
		t.Errorf("AppID: got %q, want %q (primary should win)", cfg.AppID, "primary")
	}
}

func TestLoadFromEnv_Password(t *testing.T) {
	os.Setenv("MICROSOFT_APP_PASSWORD", "my-secret")
	defer os.Unsetenv("MICROSOFT_APP_PASSWORD")
	os.Unsetenv("AZURE_CLIENT_SECRET")

	cfg := LoadFromEnv()
	if cfg.AppPassword != "my-secret" {
		t.Errorf("AppPassword: got %q, want %q", cfg.AppPassword, "my-secret")
	}
	if cfg.ClientSecret != "my-secret" {
		t.Errorf("ClientSecret: got %q, want %q", cfg.ClientSecret, "my-secret")
	}
}

func TestLoadFromEnv_TenantID(t *testing.T) {
	os.Setenv("MICROSOFT_APP_TENANT_ID", "tenant-abc")
	defer os.Unsetenv("MICROSOFT_APP_TENANT_ID")
	os.Unsetenv("AZURE_TENANT_ID")

	cfg := LoadFromEnv()
	if cfg.TenantID != "tenant-abc" {
		t.Errorf("TenantID: got %q, want %q", cfg.TenantID, "tenant-abc")
	}
}

func TestLoadFromEnv_ServiceURL(t *testing.T) {
	os.Setenv("SERVICE_URL", "https://service.example.com")
	defer os.Unsetenv("SERVICE_URL")

	cfg := LoadFromEnv()
	if cfg.ServiceURL != "https://service.example.com" {
		t.Errorf("ServiceURL: got %q, want %q", cfg.ServiceURL, "https://service.example.com")
	}
}

func TestLoadFromMap_FlatKey(t *testing.T) {
	envVars := map[string]string{
		"APP_ID": "flat-value",
	}
	result := LoadFromMap(envVars)
	if v, ok := result["APP_ID"]; !ok || v != "flat-value" {
		t.Errorf("flat key: got %v, want %q", v, "flat-value")
	}
}

func TestLoadFromMap_NestedKey(t *testing.T) {
	envVars := map[string]string{
		"CONNECTIONS__MyConn__APP_ID": "nested-app-id",
	}
	result := LoadFromMap(envVars)

	connections, ok := result["CONNECTIONS"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected CONNECTIONS to be map[string]interface{}, got %T", result["CONNECTIONS"])
	}
	myConn, ok := connections["MyConn"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected MyConn to be map[string]interface{}, got %T", connections["MyConn"])
	}
	if myConn["APP_ID"] != "nested-app-id" {
		t.Errorf("nested value: got %v, want %q", myConn["APP_ID"], "nested-app-id")
	}
}

func TestLoadFromMap_MultipleKeys(t *testing.T) {
	envVars := map[string]string{
		"A__B": "val-ab",
		"A__C": "val-ac",
	}
	result := LoadFromMap(envVars)

	aMap, ok := result["A"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected A to be map, got %T", result["A"])
	}
	if aMap["B"] != "val-ab" {
		t.Errorf("A.B: got %v, want %q", aMap["B"], "val-ab")
	}
	if aMap["C"] != "val-ac" {
		t.Errorf("A.C: got %v, want %q", aMap["C"], "val-ac")
	}
}

func TestLoadFromMap_EmptyMap(t *testing.T) {
	result := LoadFromMap(map[string]string{})
	if result == nil {
		t.Error("expected non-nil result for empty input")
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestAgentConfigStruct(t *testing.T) {
	cfg := &AgentConfig{
		ClientID:     "cid",
		ClientSecret: "csec",
		TenantID:     "tid",
		AppID:        "aid",
		AppPassword:  "pass",
		ServiceURL:   "https://svc.example.com",
		Connections: map[string]*ConnectionConfig{
			"conn1": {
				ServiceURL:    "https://conn.example.com",
				AppID:         "conn-app",
				AppPassword:   "conn-pass",
				TenantID:      "conn-tenant",
				TokenEndpoint: "https://token.example.com",
				Scopes:        []string{"scope1", "scope2"},
			},
		},
	}

	if len(cfg.Connections) != 1 {
		t.Errorf("expected 1 connection, got %d", len(cfg.Connections))
	}
	conn := cfg.Connections["conn1"]
	if conn == nil {
		t.Fatal("expected conn1 to be present")
	}
	if len(conn.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(conn.Scopes))
	}
}
