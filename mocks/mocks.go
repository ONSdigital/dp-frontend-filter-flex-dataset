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
	"[SelectAreaTypeLeadText]",
	"one = \"Select area type (cy)\"",
	"[SelectCategoriesLeadText]",
	"one = \"Select categories (cy)\"",
	"[Category]",
	"one = \"category (cy)\"",
	"other = \"categories (cy)\"",
	"[ChangeAreaTypeWarning]",
	"one = \"Saved options warning (cy)\"",
	"[DimensionsChangeWarning]",
	"one = \"Dimensions change warning (cy)\"",
	"[ImproveResultsTitle]",
	"other = \"Improve results title (cy)\"",
	"[ImproveResultsSubHeading]",
	"one = \"Improve results subheading (cy)\"",
	"[ImproveResultsList]",
	"one = \"Improve your results (cy)\"",
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
	"[SelectAreaTypeLeadText]",
	"one = \"Select area type\"",
	"[SelectCategoriesLeadText]",
	"one = \"Select categories\"",
	"[Category]",
	"one = \"category\"",
	"other = \"categories\"",
	"[ChangeAreaTypeWarning]",
	"one = \"Saved options warning\"",
	"[DimensionsChangeWarning]",
	"one = \"Dimensions change warning\"",
	"[ImproveResultsTitle]",
	"other = \"Improve results title\"",
	"[ImproveResultsSubHeading]",
	"one = \"Improve results sub heading\"",
	"[ImproveResultsList]",
	"one = \"Improve your results\"",
}

// MockAssetFunction returns mocked toml []bytes
func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
