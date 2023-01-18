package mapper

import (
	"sort"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
)

// sortCategoriesByID sorts population categories by ID - numerically if possible, with negatives listed last
func sortCategoriesByID(items []population.Category) []population.Category {
	sorted := []population.Category{}
	sorted = append(sorted, items...)

	doNumericSort := func(items []population.Category) bool {
		for _, item := range items {
			_, err := strconv.Atoi(item.ID)
			if err != nil {
				return false
			}
		}
		return true
	}

	if doNumericSort(items) {
		sort.Slice(sorted, func(i, j int) bool {
			left, _ := strconv.Atoi(sorted[i].ID)
			right, _ := strconv.Atoi(sorted[j].ID)
			if left*right < 0 {
				return right < 0
			} else {
				return left*left < right*right
			}
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ID < sorted[j].ID
		})
	}
	return sorted
}

// sortAreaTypes sorts area types by Hierarchy_Order field descending, then TotalCount
func sortAreaTypes(areaTypes []population.AreaType) []population.AreaType {
	sorted := append([]population.AreaType{}, areaTypes...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Hierarchy_Order != sorted[j].Hierarchy_Order {
			return sorted[i].Hierarchy_Order > sorted[j].Hierarchy_Order
		} else {
			return sorted[i].TotalCount < sorted[j].TotalCount
		}
	})
	return sorted
}

func filterOutWards(areaTypes []population.AreaType) []population.AreaType {
	filtered := []population.AreaType{}
	for _, areaType := range areaTypes {
		if areaType.ID != "wd" {
			filtered = append(filtered, areaType)
		}
	}
	return filtered
}
