// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package authorization

// JWT claim names used in Microsoft identity tokens.
const (
	// AudienceClaim is the "aud" claim identifying the intended recipients of the JWT.
	// See RFC 7519 §4.1.3.
	AudienceClaim = "aud"

	// IssuerClaim is the "iss" claim identifying the principal that issued the JWT.
	// See RFC 7519 §4.1.1.
	IssuerClaim = "iss"

	// AppIDClaim is the "appid" claim used in Microsoft AAD v1.0 tokens.
	AppIDClaim = "appid"

	// AppIDV2Claim ("azp") is the authorized-party claim used in AAD v2.0 tokens.
	// See OpenID Connect Core §2.
	AppIDV2Claim = "azp"

	// VersionClaim is the "ver" token-version claim used in Microsoft AAD tokens.
	VersionClaim = "ver"

	// ServiceURLClaim is the "serviceurl" claim carrying the Bot Framework service URL.
	ServiceURLClaim = "serviceurl"

	// TenantIDClaim is the "tid" tenant-ID claim used in Microsoft AAD tokens.
	TenantIDClaim = "tid"

	// KeyIDHeader is the "kid" JOSE header parameter that identifies the signing key.
	// See RFC 7515 §4.1.4.
	KeyIDHeader = "kid"
)

// Well-known Bot Framework / Agents SDK OAuth endpoints and scope values.
const (
	// AgentsSDKScope is the OAuth scope used to request Agents SDK tokens.
	AgentsSDKScope = "https://api.botframework.com"

	// AgentsSDKTokenIssuer is the token issuer for ABS-issued tokens.
	AgentsSDKTokenIssuer = "https://api.botframework.com"

	// AgentsSDKOAuthURL is the default OAuth URL for obtaining user tokens.
	AgentsSDKOAuthURL = "https://api.botframework.com"

	// PublicABSOpenIDMetadataURL is the public ABS OpenID Connect metadata endpoint.
	PublicABSOpenIDMetadataURL = "https://login.botframework.com/v1/.well-known/openidconfiguration"

	// PublicOpenIDMetadataURL is the Microsoft common v2 OpenID Connect metadata endpoint.
	PublicOpenIDMetadataURL = "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration"

	// GovABSOpenIDMetadataURL is the US-government ABS OpenID Connect metadata endpoint.
	GovABSOpenIDMetadataURL = "https://login.botframework.azure.us/v1/.well-known/openidconfiguration"

	// GovOpenIDMetadataURL is the US-government AAD OpenID Connect metadata endpoint.
	GovOpenIDMetadataURL = "https://login.microsoftonline.us/cab8a31a-1906-4287-a0d8-4eef66b95f6e/v2.0/.well-known/openid-configuration"

	// ValidTokenIssuerURLTemplateV1 is the v1 AAD issuer URL template; substitute {tenantID}.
	ValidTokenIssuerURLTemplateV1 = "https://sts.windows.net/%s/"

	// ValidTokenIssuerURLTemplateV2 is the v2 AAD issuer URL template; substitute {tenantID}.
	ValidTokenIssuerURLTemplateV2 = "https://login.microsoftonline.com/%s/v2.0"

	// EnterpriseChannelOpenIDMetadataURLFormat is the enterprise channel metadata URL; substitute {channelName}.
	EnterpriseChannelOpenIDMetadataURLFormat = "https://%s.enterprisechannel.botframework.com/v1/.well-known/openidconfiguration"

	// BotFrameworkJWKSURL is the JWKS endpoint for verifying ABS-issued tokens.
	BotFrameworkJWKSURL = "https://login.botframework.com/v1/.well-known/keys"

	// AnonymousSkillAppID is the synthetic app ID used for anonymous skill authentication.
	AnonymousSkillAppID = "anonymous"
)

// APx environment-specific OAuth scopes.
const (
	APxLocalScope      = "c16e153d-5d2b-4c21-b7f4-b05ee5d516f1/.default"
	APxDevScope        = "0d94caae-b412-4943-8a68-83135ad6d35f/.default"
	APxProductionScope = "5a807f24-c9de-44ee-a3a7-329e88a00ffc/.default"
	APxGCCScope        = "c9475445-9789-4fef-9ec5-cde4a9bcd446/.default"
	APxGCCHScope       = "6f669b9e-7701-4e2b-b624-82c9207fde26/.default"
	APxDoDScope        = "0a069c81-8c7c-4712-886b-9c542d673ffb/.default"
	APxGallatin        = "bd004c8e-5acf-4c48-8570-4e7d46b2f63b/.default"
)

// ChannelAdapter turn-state keys used to pass values through the pipeline.
const (
	AgentIdentityKey          = "AgentIdentity"
	OAuthScopeKey             = "Microsoft.Agents.Builder.ChannelAdapter.OAuthScope"
	InvokeResponseKey         = "ChannelAdapter.InvokeResponse"
	ConnectorFactoryKey       = "ConnectorFactory"
	UserTokenClientKey        = "UserTokenClient"
	AgentCallbackHandlerKey   = "AgentCallbackHandler"
	ChannelServiceFactoryKey  = "ChannelServiceClientFactory"
)
