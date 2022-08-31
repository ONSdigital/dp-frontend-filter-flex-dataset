package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
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

// GetCoverage handler
func GetCoverage(rc RenderClient, fc FilterClient, pc PopulationClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, fc, pc, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, pc PopulationClient, lang, accessToken, collectionID string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	c := req.URL.Query().Get("c")
	q := req.URL.Query().Get("q")
	pq := req.URL.Query().Get("pq")
	p := req.URL.Query().Get("p")
	isNameSearch := strings.Contains(req.URL.RawQuery, "q=")
	isParentSearch := strings.Contains(req.URL.RawQuery, "p=")
	isValidationError, _ := strconv.ParseBool(req.URL.Query().Get("error"))
	var filterJob *filter.GetFilterResponse
	var filterDims filter.Dimensions
	var parents population.GetAreaTypeParentsResponse
	var opts filter.DimensionOptions
	var areas population.GetAreasResponse
	var fErr, dErr, pErr, oErr, nsErr, psErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		filterJob, fErr = fc.GetFilter(ctx, filter.GetFilterInput{
			FilterID: filterID,
			AuthHeaders: filter.AuthHeaders{
				UserAuthToken: accessToken,
				CollectionID:  collectionID,
			},
		})
	}()
	go func() {
		defer wg.Done()
		filterDims, _, dErr = fc.GetDimensions(ctx, accessToken, "", collectionID, filterID, &filter.QueryParams{Limit: 500})
	}()
	wg.Wait()

	if fErr != nil {
		log.Error(ctx, "failed to get filter", fErr, log.Data{"filter_id": filterID})
		setStatusCode(req, w, fErr)
		return
	}
	if dErr != nil {
		log.Error(ctx, "failed to get dimensions", dErr, log.Data{"filter_id": filterID})
		setStatusCode(req, w, dErr)
		return
	}

	var geogLabel, geogID, dimension, parent string
	for _, dim := range filterDims.Items {
		// Needed to determine whether dimension is_area_type
		// Only one dimension will be is_area_type=true
		filterDimension, _, err := fc.GetDimension(ctx, accessToken, "", collectionID, filterID, dim.Name)
		if err != nil {
			log.Error(ctx, "failed to get dimension", err, log.Data{"dimension_name": dim.Name})
			setStatusCode(req, w, err)
			return
		}
		if *filterDimension.IsAreaType {
			geogLabel = filterDimension.Label
			geogID = filterDimension.ID
			dimension = filterDimension.Name
			parent = filterDimension.FilterByParent
			break
		}
	}

	var hasFilterByParent bool
	if parent != "" {
		hasFilterByParent = true
	}

	wg.Add(4)
	go func() {
		defer wg.Done()
		parents, pErr = pc.GetAreaTypeParents(ctx, population.GetAreaTypeParentsInput{
			UserAuthToken: accessToken,
			DatasetID:     filterJob.PopulationType,
			AreaTypeID:    geogID,
		})
	}()
	go func() {
		defer wg.Done()
		opts, _, oErr = fc.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dimension, &filter.QueryParams{})
	}()
	go func() {
		defer wg.Done()
		if isNameSearch && q != "" {
			areas, nsErr = getAreas(pc, ctx, accessToken, filterJob.PopulationType, geogID, q)
		}
	}()
	go func() {
		defer wg.Done()
		if isParentSearch && pq != "" {
			areas, psErr = getAreas(pc, ctx, accessToken, filterJob.PopulationType, p, pq)
		}
	}()
	wg.Wait()

	if pErr != nil {
		log.Error(ctx, "failed to get parents", pErr, log.Data{
			"dataset_id":   geogID,
			"area_type_id": geogLabel,
		})
		setStatusCode(req, w, pErr)
		return
	}
	if oErr != nil {
		log.Error(ctx, "failed to get dimension options", oErr, log.Data{"dimension_name": dimension})
		setStatusCode(req, w, oErr)
		return
	}
	if nsErr != nil {
		log.Error(ctx, "failed to get areas in name search", nsErr, log.Data{
			"population_type": filterJob.PopulationType,
			"area":            geogID,
			"query":           q,
		})
		setStatusCode(req, w, nsErr)
		return
	}
	if psErr != nil {
		log.Error(ctx, "failed to get areas in parent search", psErr, log.Data{
			"population_type": filterJob.PopulationType,
			"area":            p,
			"query":           pq,
		})
		setStatusCode(req, w, psErr)
		return
	}

	options := []model.SelectableElement{}
	var areaType string
	if hasFilterByParent {
		areaType = parent
	} else {
		areaType = geogID
	}
	for _, opt := range opts.Items {
		var option model.SelectableElement

		area, err := pc.GetArea(ctx, population.GetAreaInput{
			UserAuthToken:  accessToken,
			PopulationType: filterJob.PopulationType,
			AreaType:       areaType,
			Area:           opt.Option,
		})
		if err != nil {
			log.Error(ctx, "failed to get area", err, log.Data{
				"population": filterJob.PopulationType,
				"area type":  geogID,
				"ID":         opt.Option,
			})
			setStatusCode(req, w, err)
			return
		}
		option.Value = opt.Option
		option.Text = area.Area.Label
		log.Warn(ctx, "AreaLabelXX", log.Data{
			"area": area.Area,
		})

		options = append(options, option)
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geogLabel, q, pq, p, c, dimension, geogID, areas, options, parents, hasFilterByParent, isValidationError)
	rc.BuildPage(w, m, "coverage")
}

// getAreas is a helper function that returns the GetAreasResponse or an error
func getAreas(pc PopulationClient, ctx context.Context, accessToken, popType, areaTypeID, query string) (population.GetAreasResponse, error) {
	areas, err := pc.GetAreas(ctx, population.GetAreasInput{
		UserAuthToken: accessToken,
		DatasetID:     popType,
		AreaTypeID:    areaTypeID,
		Text:          url.QueryEscape(strings.TrimSpace(query)),
	})
	return areas, err
}
