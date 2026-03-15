// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package app

import (
	"encoding/json"
	"fmt"
)

// Query extracts typed data from an activity's Value field.
// It is used to parse structured data from invoke activities,
// messaging extension queries, and adaptive card submissions.
type Query[T any] struct {
	// Data holds the parsed value.
	Data T
}

// ParseQuery extracts and parses the activity value into a Query[T].
// value should be the activity.Value field (interface{} or map).
func ParseQuery[T any](value interface{}) (*Query[T], error) {
	if value == nil {
		return &Query[T]{}, nil
	}

	// Marshal to JSON then unmarshal into T
	b, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("query: marshaling value: %w", err)
	}

	var data T
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("query: unmarshaling into type: %w", err)
	}

	return &Query[T]{Data: data}, nil
}
