// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package connector

import "context"

// ChannelServiceClientFactory creates connector clients based on channel context.
type ChannelServiceClientFactory interface {
	// CreateConnectorClient creates a connector client for the given service URL.
	CreateConnectorClient(serviceURL string) (*ConnectorClient, error)
	// CreateUserTokenClient creates a user token client.
	CreateUserTokenClient(serviceURL string) (*UserTokenClient, error)
}

// contextTokenProvider wraps the factory's token provider interface to match
// the TokenProvider interface used by ConnectorClient.
type contextTokenProvider struct {
	inner interface {
		GetAccessToken(ctx context.Context, resource string) (string, error)
	}
}

func (p *contextTokenProvider) GetAccessToken(ctx context.Context, resource string) (string, error) {
	return p.inner.GetAccessToken(ctx, resource)
}

// RestChannelServiceClientFactory is the default factory implementation.
type RestChannelServiceClientFactory struct {
	tokenProvider interface {
		GetAccessToken(ctx context.Context, resource string) (string, error)
	}
}

// NewRestChannelServiceClientFactory creates a new RestChannelServiceClientFactory.
func NewRestChannelServiceClientFactory(tokenProvider interface {
	GetAccessToken(ctx context.Context, resource string) (string, error)
}) *RestChannelServiceClientFactory {
	return &RestChannelServiceClientFactory{
		tokenProvider: tokenProvider,
	}
}

// CreateConnectorClient creates a ConnectorClient for the given service URL.
func (f *RestChannelServiceClientFactory) CreateConnectorClient(serviceURL string) (*ConnectorClient, error) {
	var tp TokenProvider
	if f.tokenProvider != nil {
		tp = &contextTokenProvider{inner: f.tokenProvider}
	}
	return NewConnectorClient(serviceURL, tp), nil
}

// CreateUserTokenClient creates a UserTokenClient for the given service URL.
func (f *RestChannelServiceClientFactory) CreateUserTokenClient(serviceURL string) (*UserTokenClient, error) {
	return NewUserTokenClient(serviceURL), nil
}
