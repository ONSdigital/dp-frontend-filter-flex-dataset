package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// ChangeDimensions Handler
func ChangeDimensions(rc RenderClient, fc FilterClient, dc DatasetClient, pc PopulationClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		changeDimensions(w, req, rc, fc, dc, pc, accessToken, collectionID, lang)
	})
}

func changeDimensions(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, dc DatasetClient, pc PopulationClient, accessToken, collectionID, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	fid := vars["filterID"]
	var fErr, imErr, pErr error
	var pDims population.GetDimensionsResponse
	var dims []model.FilterDimension
	var popType string
	var isMultivariate bool

	// get filter dimensions
	fDims, _, err := fc.GetDimensions(ctx, accessToken, "", collectionID, fid, &filter.QueryParams{Limit: 500})
	if err != nil {
		log.Error(ctx, "failed to get dimensions", err, log.Data{"filter_id": fid})
		setStatusCode(req, w, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var fj *filter.GetFilterResponse
		fj, fErr = fc.GetFilter(ctx, filter.GetFilterInput{
			FilterID: fid,
			AuthHeaders: filter.AuthHeaders{
				UserAuthToken: accessToken,
				CollectionID:  collectionID,
			},
		})
		popType = fj.PopulationType

		// check dataset is multivariate
		isMultivariate, imErr = isMultivariateDataset(dc, ctx, accessToken, collectionID, fj.Dataset.DatasetID)
		if !isMultivariate && imErr == nil {
			http.Redirect(w, req, fmt.Sprintf("/filters/%s/dimensions", fid), http.StatusMovedPermanently)
			return
		}

		// get available population dimensions
		pDims, pErr = pc.GetDimensions(ctx, population.GetDimensionsInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PaginationParams: population.PaginationParams{
				Limit: 1000,
			},
			PopulationType: popType,
		})
	}()

	dimErrs := make([]error, len(fDims.Items))
	go func() {
		defer wg.Done()
		for i, dim := range fDims.Items {
			// Needed to determine whether dimension is_area_type
			fDim, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, fid, dim.Name)
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

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetChangeDimensions(req, basePage, lang, fid, dims, pDims)
	rc.BuildPage(w, m, "dimensions")
}

// isMultivariateDataset determines whether the given filter record is based on a multivariate dataset type
func isMultivariateDataset(dc DatasetClient, ctx context.Context, accessToken, collectionID, did string) (bool, error) {
	d, err := dc.Get(ctx, accessToken, "", collectionID, did)
	if err != nil {
		return false, fmt.Errorf("failed to get dataset: %w", err)
	}

	if strings.Contains(d.Type, "multivariate") {
		return true, nil
	}
	return false, nil
}
