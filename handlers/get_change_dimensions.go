package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetChangeDimensions Handler
func (f *FilterFlex) GetChangeDimensions() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getChangeDimensions(w, req, f, accessToken, collectionID, lang)
	})
}

func getChangeDimensions(w http.ResponseWriter, req *http.Request, f *FilterFlex, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	fid := vars["filterID"]
	q := req.URL.Query().Get("q")
	isSearch := strings.Contains(req.URL.RawQuery, "q=")
	form := req.URL.Query().Get("f")
	var fErr, imErr, pErr, prErr, zErr error
	var pDims, pResults population.GetDimensionsResponse
	var dims []model.FilterDimension
	var eb zebedee.EmergencyBanner
	var popType, serviceMsg string
	var isMultivariate bool

	// get filter dimensions
	fDims, _, err := f.FilterClient.GetDimensions(ctx, accessToken, "", collectionID, fid, &filter.QueryParams{Limit: 500})
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": fid})
		setStatusCode(req, w, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		eb, serviceMsg, zErr = getZebContent(ctx, f.ZebedeeClient, accessToken, collectionID, lang)
	}()

	go func() {
		defer wg.Done()
		var fj *filter.GetFilterResponse
		fj, fErr = f.FilterClient.GetFilter(ctx, filter.GetFilterInput{
			FilterID: fid,
			AuthHeaders: filter.AuthHeaders{
				UserAuthToken: accessToken,
				CollectionID:  collectionID,
			},
		})
		popType = fj.PopulationType

		// check dataset is multivariate
		isMultivariate, imErr = isMultivariateDataset(ctx, f.DatasetClient, accessToken, collectionID, fj.Dataset.DatasetID)
		if !isMultivariate && imErr == nil {
			http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", fid), http.StatusMovedPermanently)
			return
		}

		// get available population dimensions
		pDims, pErr = f.PopulationClient.GetDimensions(ctx, population.GetDimensionsInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PaginationParams: population.PaginationParams{
				Limit: 1000,
			},
			PopulationType: popType,
		})

		if isSearch && q != "" {
			pResults, prErr = f.PopulationClient.GetDimensions(ctx, population.GetDimensionsInput{
				AuthTokens: population.AuthTokens{
					UserAuthToken: accessToken,
				},
				PaginationParams: population.PaginationParams{
					Limit: 1000,
				},
				PopulationType: popType,
				SearchString:   url.QueryEscape(strings.TrimSpace(q)),
			})
		}
	}()

	dimErrs := make([]error, len(fDims.Items))
	go func() {
		defer wg.Done()
		for i, dim := range fDims.Items {
			// Needed to determine whether dimension is_area_type
			fDim, _, err := f.FilterClient.GetDimension(ctx, accessToken, "", collectionID, fid, dim.Name)
			if err != nil {
				log.Error(ctx, "failed to get dimension", err, log.Data{
					"dimension_name": dim.Name,
				})
				dimErrs[i] = err
			}
			dim.IsAreaType = fDim.IsAreaType
			dims = append(dims, model.FilterDimension{
				Dimension: dim,
			})
		}
	}()
	wg.Wait()

	// error handling from waitgroup
	// log zebedee error but don't set a server error
	if zErr != nil {
		log.Error(ctx, "unable to get homepage content", zErr, log.Data{"homepage_content": zErr})
	}
	if fErr != nil {
		log.Error(ctx, "failed to get filter", fErr, log.Data{
			"filter_id": fid,
		})
		setStatusCode(req, w, fErr)
		return
	}
	if imErr != nil {
		log.Error(ctx, "failed to determine if dataset type is multivariate", imErr, log.Data{
			"filter_id": fid,
		})
		setStatusCode(req, w, imErr)
		return
	}
	if pErr != nil {
		log.Error(ctx, "failed to get population dimensions", pErr, log.Data{
			"population_type": popType,
		})
		setStatusCode(req, w, pErr)
		return
	}
	if prErr != nil {
		log.Error(ctx, "failed to get population dimensions from query", prErr, log.Data{
			"population_type": popType,
			"query":           q,
		})
		setStatusCode(req, w, prErr)
		return
	}
	var hasErrs bool
	for _, err := range dimErrs {
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{
				"filter_id": fid,
			})
			hasErrs = true
		}
	}
	if hasErrs {
		setStatusCode(req, w, err)
		return
	}

	if !isMultivariate {
		http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", fid), http.StatusMovedPermanently)
		return
	}

	basePage := f.Render.NewBasePageModel()
	m := mapper.NewMapper(req, basePage, eb, lang, serviceMsg, fid)
	dimensions := m.CreateGetChangeDimensions(q, form, dims, pDims, pResults)
	f.Render.BuildPage(w, dimensions, "dimensions")
}
