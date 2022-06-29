package mocks

import "strings"

var cyLocale = []string{
	"[Back]",
	"one = \"Back\"",
}

var enLocale = []string{
	"[Back]",
	"one = \"Back\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
