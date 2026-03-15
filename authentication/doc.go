// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package authentication provides MSAL-based OAuth authentication for the
// Microsoft 365 Agents SDK for Go.
//
// It implements the AccessTokenProvider interface from hosting/core,
// supporting client secret, certificate, and managed identity authentication.
//
// Basic usage with client secret:
//
//	auth, err := authentication.NewMsalAuth(authentication.Config{
//	    TenantID:     os.Getenv("TENANT_ID"),
//	    ClientID:     os.Getenv("CLIENT_ID"),
//	    ClientSecret: os.Getenv("CLIENT_SECRET"),
//	})
//	manager := authentication.NewConnectionManager(auth)
package authentication
