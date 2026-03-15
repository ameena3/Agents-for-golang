// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package activity

// EntityType constants for known entity types.
const (
	EntityTypeMention     = "mention"
	EntityTypePlace       = "Place"
	EntityTypeGeoCoords   = "GeoCoordinates"
	EntityTypeProductInfo = "ProductInfo"
)

// Entity represents metadata pertaining to an activity.
type Entity struct {
	// Type is the type of this entity (RFC 3987 IRI).
	Type string `json:"type,omitempty"`
	// AdditionalProperties holds extra entity properties.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Mention represents a mention entity (entity type: "mention").
type Mention struct {
	// Type is always "mention".
	Type string `json:"type,omitempty"`
	// Mentioned is the mentioned user.
	Mentioned *ChannelAccount `json:"mentioned,omitempty"`
	// Text is the sub-text which represents the mention.
	Text string `json:"text,omitempty"`
}

// GeoCoordinates represents a geographical location.
type GeoCoordinates struct {
	// Type is always "GeoCoordinates".
	Type string `json:"type,omitempty"`
	// Elevation is the elevation of the location.
	Elevation float64 `json:"elevation,omitempty"`
	// Latitude is the latitude of the location.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude is the longitude of the location.
	Longitude float64 `json:"longitude,omitempty"`
	// Name is the name of this location.
	Name string `json:"name,omitempty"`
}

// Place represents a place entity.
type Place struct {
	// Type is always "Place".
	Type string `json:"type,omitempty"`
	// Address is the address of the place.
	Address interface{} `json:"address,omitempty"`
	// Geo is the geographical coordinates of the place.
	Geo *GeoCoordinates `json:"geo,omitempty"`
	// HasMap is the URL to a map of the place.
	HasMap interface{} `json:"hasMap,omitempty"`
	// Name is the name of the place.
	Name string `json:"name,omitempty"`
}

// ClientCitationImage contains information about a citation's icon.
type ClientCitationImage struct {
	// Type is the schema type, typically "ImageObject".
	Type string `json:"type,omitempty"`
	// Name is the name of the image.
	Name string `json:"name,omitempty"`
}

// SensitivityPattern contains pattern information for sensitivity usage info.
type SensitivityPattern struct {
	// Type is the schema type, typically "DefinedTerm".
	Type string `json:"type,omitempty"`
	// InDefinedTermSet is the term set this pattern belongs to.
	InDefinedTermSet string `json:"inDefinedTermSet,omitempty"`
	// Name is the name of the pattern.
	Name string `json:"name,omitempty"`
	// TermCode is the code for this term.
	TermCode string `json:"termCode,omitempty"`
}

// SensitivityUsageInfo contains sensitivity usage information for AI content.
type SensitivityUsageInfo struct {
	// Type is the schema type.
	Type string `json:"type,omitempty"`
	// SchemaType is the schema type identifier.
	SchemaType string `json:"@type,omitempty"`
	// Description is an optional description.
	Description string `json:"description,omitempty"`
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Position is the optional position.
	Position *int `json:"position,omitempty"`
	// Pattern is the optional sensitivity pattern.
	Pattern *SensitivityPattern `json:"pattern,omitempty"`
}

// ClientCitationAppearance contains appearance information for a client citation.
type ClientCitationAppearance struct {
	// Type is the schema type, typically "DigitalDocument".
	Type string `json:"type,omitempty"`
	// Name is the name of the citation.
	Name string `json:"name,omitempty"`
	// Text is optional text content.
	Text string `json:"text,omitempty"`
	// URL is the optional URL for the citation.
	URL string `json:"url,omitempty"`
	// Abstract is a short abstract of the citation.
	Abstract string `json:"abstract,omitempty"`
	// EncodingFormat is the optional encoding format.
	EncodingFormat string `json:"encodingFormat,omitempty"`
	// Image is the optional citation image.
	Image *ClientCitationImage `json:"image,omitempty"`
	// Keywords are optional keywords for the citation.
	Keywords []string `json:"keywords,omitempty"`
	// UsageInfo is optional sensitivity usage info.
	UsageInfo *SensitivityUsageInfo `json:"usageInfo,omitempty"`
}

// ClientCitation represents a Teams client citation to include in a message.
type ClientCitation struct {
	// Type is always "Claim".
	Type string `json:"type,omitempty"`
	// Position is the position of the citation.
	Position int `json:"position,omitempty"`
	// Appearance is the appearance information for the citation.
	Appearance *ClientCitationAppearance `json:"appearance,omitempty"`
}

// AIEntity represents an entity that indicates AI-generated content.
type AIEntity struct {
	// Type is the schema type.
	Type string `json:"type,omitempty"`
	// SchemaType is the schema type value.
	SchemaType string `json:"@type,omitempty"`
	// Context is the schema context.
	Context string `json:"@context,omitempty"`
	// ID is the entity ID.
	ID string `json:"@id,omitempty"`
	// AdditionalType contains additional type values.
	AdditionalType []string `json:"additionalType,omitempty"`
	// Citation contains the citations for this AI entity.
	Citation []*ClientCitation `json:"citation,omitempty"`
	// UsageInfo contains the sensitivity usage info.
	UsageInfo *SensitivityUsageInfo `json:"usageInfo,omitempty"`
}
