package mapper

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSortCategoriesByID(t *testing.T) {
	Convey("Population categories are sorted", t, func() {
		getIDList := func(items []population.Category) []string {
			results := []string{}
			for _, item := range items {
				results = append(results, item.ID)
			}
			return results
		}

		Convey("given non-numeric options", func() {
			nonNumeric := []population.Category{
				{
					ID:    "dim_2",
					Label: "option 2",
				},
				{
					ID:    "dim_1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(nonNumeric)

				Convey("then options are sorted alphabetically by ID", func() {
					actual := getIDList(sorted)
					expected := []string{"dim_1", "dim_2"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given simple numeric options", func() {
			numeric := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2",
					Label: "option 2",
				}, {
					ID:    "1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(numeric)

				Convey("then options are sorted numerically", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "2", "10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given numeric options with negatives", func() {
			numericWithNegatives := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2",
					Label: "option 2",
				}, {
					ID:    "-1",
					Label: "option -1",
				},
				{
					ID:    "1",
					Label: "option 1",
				}, {
					ID:    "-10",
					Label: "option -10",
				},
			}

			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(numericWithNegatives)

				Convey("then options are sorted numerically with negatives at the end", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "2", "10", "-1", "-10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given mixed numeric and non-numeric options", func() {
			alphanumeric := []population.Category{
				{
					ID:    "10",
					Label: "option 10",
				}, {
					ID:    "2nd Option",
					Label: "option 2",
				}, {
					ID:    "1",
					Label: "option 1",
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortCategoriesByID(alphanumeric)

				Convey("then options are sorted alphanumerically", func() {
					actual := getIDList(sorted)
					expected := []string{"1", "10", "2nd Option"}
					So(actual, ShouldResemble, expected)
				})
			})
		})
	})
}

func TestSortAreaTypes(t *testing.T) {
	Convey("Population AreaTypes are sorted", t, func() {
		getIDList := func(items []population.AreaType) []string {
			results := []string{}
			for _, item := range items {
				results = append(results, item.ID)
			}
			return results
		}

		Convey("given valid heirarchy_order fields", func() {
			areaTypes := []population.AreaType{
				{
					ID:              "H_100",
					TotalCount:      3,
					Hierarchy_Order: 100,
				},
				{
					ID:              "H_10",
					TotalCount:      2,
					Hierarchy_Order: 10,
				},
				{
					ID:              "H_20",
					TotalCount:      1,
					Hierarchy_Order: 20,
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortAreaTypes(areaTypes)

				Convey("then AreaTypes are sorted descending numerically by Hierarchy_Order", func() {
					actual := getIDList(sorted)
					expected := []string{"H_100", "H_20", "H_10"}
					So(actual, ShouldResemble, expected)
				})
			})
		})

		Convey("given empty heirarchy_order fields", func() {
			areaTypes := []population.AreaType{
				{
					ID:         "H_100",
					TotalCount: 3,
				},
				{
					ID:         "H_10",
					TotalCount: 2,
				},
				{
					ID:         "H_20",
					TotalCount: 1,
				},
			}
			Convey("when they are sorted", func() {
				sorted := sortAreaTypes(areaTypes)

				Convey("then AreaTypes are sorted ascending numerically by TotalCount", func() {
					actual := getIDList(sorted)
					expected := []string{"H_20", "H_10", "H_100"}
					So(actual, ShouldResemble, expected)
				})
			})
		})
	})
}
