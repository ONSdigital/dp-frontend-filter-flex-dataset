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
	"[AreasAddedTitle]",
	"one = \"Area added (cy)\"",
	"other = \"Areas added (cy)\"",
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
	"[AreasAddedTitle]",
	"one = \"Area added\"",
	"other = \"Areas added\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
