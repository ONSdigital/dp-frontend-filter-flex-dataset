package helpers

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ONSdigital/dp-renderer/helper"
	"github.com/ONSdigital/log.go/v2/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

/* Pluralise performs a toml file lookup to return the pluralised lookup value from a given key.
This function is intended not to panic if the toml file lookup fails.
req is the *http.Request, required for logging.
key is the key to lookup in the toml file.
lang is the language.
keyPrefix is an optional prefix to the key string.
plural is the plural int required for toml file lookup.
*/
func Pluralise(req *http.Request, key, lang, keyPrefix string, plural int) string {
	ctx := req.Context()
	// Convert input to a string that can be used as a lookup in toml files
	// String will be title case with given prefix e.g. PrefixInputString
	log.Info(ctx, "converting input string", log.Data{
		"input": key,
	})
	str := cases.Title(language.English).String(key)
	str = strings.ReplaceAll(str, " ", "")
	str = keyPrefix + str
	// Localise will not return an error but will panic and return an empty string.
	// The given input may not be stored in the toml files, so risk of panic is high.
	// Therefore log panic and handle returned empty string in calling function
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx, "recovered from panic", log.Data{
				"lookup_not_found": str,
			})
		}
	}()
	log.Info(ctx, "performing toml lookup", log.Data{
		"lookup": str,
	})
	return helper.Localise(str, lang, plural)
}

// TrimCategoryValue trims _[0-9] from the given string and returns the result
func TrimCategoryValue(s string) string {
	rx := regexp.MustCompile(`(_[\d])\w+`)
	return rx.ReplaceAllString(s, "")
}
