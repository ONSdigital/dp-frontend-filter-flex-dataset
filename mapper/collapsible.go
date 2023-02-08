package mapper

import (
	"fmt"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	coreModel "github.com/ONSdigital/dp-renderer/model"
)

type Link struct {
	Uri  string
	Text string
}

func mapImproveResultsCollapsible(dims []model.Dimension) (areaTypeUri string, linksItem string) {
	var dimsLinks []Link
	for _, dim := range dims {
		if dim.IsAreaType {
			areaTypeUri = dim.URI
		} else if dim.Name != "" {
			dimsLinks = append(dimsLinks, Link{
				Uri:  dim.URI,
				Text: dim.Name,
			})
		}
	}

	return areaTypeUri, buildLinksString(dimsLinks)
}

func buildLinksString(dimsLinks []Link) (linkStr string) {
	var penultimateItem = len(dimsLinks) - 2
	for i, link := range dimsLinks {
		switch {
		case i < penultimateItem:
			linkStr += fmt.Sprintf("<a href=\"%s\">%s</a>, ", link.Uri, link.Text)
		case i == penultimateItem:
			linkStr += fmt.Sprintf("<a href=\"%s\">%s</a> or ", link.Uri, link.Text)
		default:
			linkStr += fmt.Sprintf("<a href=\"%s\">%s</a>", link.Uri, link.Text)
		}
	}
	return linkStr
}

func mapDescriptionsCollapsible(dimDescriptions population.GetDimensionsResponse, dims []model.Dimension) []coreModel.CollapsibleItem {
	var collapsibleContentItems []coreModel.CollapsibleItem
	var areaItem coreModel.CollapsibleItem

	for _, dim := range dims {
		for _, dimDescription := range dimDescriptions.Dimensions {
			if dim.ID == dimDescription.ID && !dim.IsAreaType {
				collapsibleContentItems = append(collapsibleContentItems, coreModel.CollapsibleItem{
					Subheading: cleanDimensionLabel(dimDescription.Label),
					Content:    strings.Split(dimDescription.Description, "\n"),
				})
			} else if dim.ID == dimDescription.ID && dim.IsAreaType {
				areaItem.Subheading = cleanDimensionLabel(dimDescription.Label)
				areaItem.Content = strings.Split(dimDescription.Description, "\n")
			}
		}
	}

	collapsibleContentItems = append([]coreModel.CollapsibleItem{
		{
			Subheading: areaTypeTitle,
			SafeHTML: coreModel.Localisation{
				LocaleKey: "VariableInfoAreaType",
				Plural:    1,
			},
		},
		areaItem,
		{
			Subheading: coverageTitle,
			SafeHTML: coreModel.Localisation{
				LocaleKey: "VariableInfoCoverage",
				Plural:    1,
			},
		},
	}, collapsibleContentItems...)

	return collapsibleContentItems
}
