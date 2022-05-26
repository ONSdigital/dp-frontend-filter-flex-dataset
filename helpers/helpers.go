package helpers

import (
	"net/url"
)

// HasStringInSlice checks for a string within in a string array
func HasStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// PersistExistingParams persists existing query string values and ignores a given value
func PersistExistingParams(values []string, key, ignoreValue string, q url.Values) {
	if len(values) > 0 {
		for _, value := range values {
			if value != ignoreValue {
				q.Add(key, value)
			}
		}
	}
}

// ToBoolPtr converts a boolean to a pointer
func ToBoolPtr(val bool) *bool {
	return &val
}
