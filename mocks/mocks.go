package mocks

import "strings"

var cyLocale = []string{
	"[Back]",
	"one = \"Back\"",
	"[AreaTypeCountry]",
	"one = \"Country\"",
	"other = \"Countries\"",
	"[Test]",
	"one = \"Test (cy)\"",
	"other = \"Tests (cy)\"",
	"[CoverageSelectDefault]",
	"one = \"Select (cy)\"",
}

var enLocale = []string{
	"[Back]",
	"one = \"Back\"",
	"[AreaTypeCountry]",
	"one = \"Country\"",
	"other = \"Countries\"",
	"[Test]",
	"one = \"Test\"",
	"other = \"Tests\"",
	"[CoverageSelectDefault]",
	"one = \"Select\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
