// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package client

// PowerPlatformCloud represents the Power Platform cloud environment.
type PowerPlatformCloud string

const (
	// PowerPlatformCloudPublic is the public commercial cloud.
	PowerPlatformCloudPublic PowerPlatformCloud = "Public"
	// PowerPlatformCloudGCC is the US Government Community Cloud.
	PowerPlatformCloudGCC PowerPlatformCloud = "GCC"
	// PowerPlatformCloudGCCHigh is the US Government Community Cloud High.
	PowerPlatformCloudGCCHigh PowerPlatformCloud = "GCCHigh"
	// PowerPlatformCloudDoD is the US Department of Defense cloud.
	PowerPlatformCloudDoD PowerPlatformCloud = "DoD"
	// PowerPlatformCloudGallatin is the Microsoft China (Gallatin) cloud.
	PowerPlatformCloudGallatin PowerPlatformCloud = "Gallatin"
)

// AgentType represents the type of Copilot Studio agent.
type AgentType string

const (
	// AgentTypePublished is a published agent.
	AgentTypePublished AgentType = "Published"
	// AgentTypePreview is a preview/test agent.
	AgentTypePreview AgentType = "Preview"
)

// PowerPlatformEnvironment holds environment details.
type PowerPlatformEnvironment struct {
	EnvironmentID string
	Cloud         PowerPlatformCloud
}

// GetDirectToEngineURL returns the Direct-to-Engine base URL for this environment.
func (e *PowerPlatformEnvironment) GetDirectToEngineURL() string {
	switch e.Cloud {
	case PowerPlatformCloudGCC:
		return "https://gcc.api.powerplatform.us"
	case PowerPlatformCloudGCCHigh:
		return "https://high.api.powerplatform.us"
	case PowerPlatformCloudDoD:
		return "https://dod.api.powerplatform.us"
	case PowerPlatformCloudGallatin:
		return "https://api.powerplatform.partner.microsoftonline.cn"
	default:
		return "https://api.powerplatform.com"
	}
}
