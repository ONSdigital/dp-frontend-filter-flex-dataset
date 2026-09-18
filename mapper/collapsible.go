package mapper

import (
	"fmt"
	"strings"

	core "github.com/ONSdigital/dis-design-system-go/v2/model"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
)

type Link struct {
	Uri  string
	Text string
}

func mapImproveResultsCollapsible(dims []model.Dimension) (areaTypeUri, linksItem string) {
	var dimsLinks []Link
	for i := range dims {
		dim := &dims[i]
		if dim.IsGeography {
			areaTypeUri = dim.URI
		} else if dim.Name != "" && dim.HasChange {
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
			linkStr += fmt.Sprintf("<a href=%q>%s</a>, ", link.Uri, link.Text)
		case i == penultimateItem:
			linkStr += fmt.Sprintf("<a href=%q>%s</a> or ", link.Uri, link.Text)
		default:
			linkStr += fmt.Sprintf("<a href=%q>%s</a>", link.Uri, link.Text)
		}
	}
	return linkStr
}

func mapDescriptionsCollapsible(dimDescriptions population.GetDimensionsResponse, dims []model.Dimension) []core.CollapsibleItem {
	var collapsibleContentItems []core.CollapsibleItem
	var areaItem core.CollapsibleItem

	for i := range dims {
		dim := &dims[i]
		for _, dimDescription := range dimDescriptions.Dimensions {
			if dim.ID == dimDescription.ID && !dim.IsGeography {
				collapsibleContentItems = append(collapsibleContentItems, core.CollapsibleItem{
					Subheading: cleanDimensionLabel(dimDescription.Label),
					Content:    strings.Split(dimDescription.Description, "\n"),
				})
			} else if dim.ID == dimDescription.ID && dim.IsGeography {
				areaItem.Subheading = cleanDimensionLabel(dimDescription.Label)
				areaItem.Content = strings.Split(dimDescription.Description, "\n")
			}
		}
	}

	collapsibleContentItems = append([]core.CollapsibleItem{
		{
			Subheading: areaTypeTitle,
			SafeHTML: core.Localisation{
				LocaleKey: "VariableInfoAreaType",
				Plural:    1,
			},
		},
		areaItem,
		{
			Subheading: coverageTitle,
			SafeHTML: core.Localisation{
				LocaleKey: "VariableInfoCoverage",
				Plural:    1,
			},
		},
	}, collapsibleContentItems...)

	return collapsibleContentItems
}
