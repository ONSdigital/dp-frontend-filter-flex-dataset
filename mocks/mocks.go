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
	"[DimensionsSearchLabel]",
	"one = \"Dimensions search label (cy)\"",
	"[CoverageSearchLabel]",
	"one = \"Coverage search label (cy)\"",
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
	"[DimensionsSearchLabel]",
	"one = \"Dimensions search label\"",
	"[CoverageSearchLabel]",
	"one = \"Coverage search label\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
