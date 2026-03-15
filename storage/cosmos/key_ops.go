// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cosmos

import "strings"

// SanitizeKey makes a key safe for use as a Cosmos DB item ID.
// Cosmos DB item IDs cannot contain /, \, ?, or #.
func SanitizeKey(key string) string {
	r := strings.NewReplacer(
		"/", "-FSLASH-",
		"\\", "-BSLASH-",
		"?", "-QMARK-",
		"#", "-HASH-",
	)
	return r.Replace(key)
}
