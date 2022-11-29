package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/model"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetCoverage handler
func GetCoverage(rc RenderClient, fc FilterClient, pc PopulationClient, dc DatasetClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		getCoverage(w, req, rc, fc, pc, dc, lang, accessToken, collectionID)
	})
}

func getCoverage(w http.ResponseWriter, req *http.Request, rc RenderClient, fc FilterClient, pc PopulationClient, dc DatasetClient, lang, accessToken, collectionID string) {
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
	var datasetDetails dataset.DatasetDetails
	var releaseDate string

	var fErr, dErr, pErr, oErr, nsErr, psErr, dsErr, rdErr error

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

	wg.Add(6)
	go func() {
		defer wg.Done()
		parents, pErr = pc.GetAreaTypeParents(ctx, population.GetAreaTypeParentsInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PaginationParams: population.PaginationParams{
				Limit: 1000,
			},
			PopulationType: filterJob.PopulationType,
			AreaTypeID:     geogID,
		})
	}()
	go func() {
		defer wg.Done()
		datasetDetails, dsErr = dc.Get(ctx, accessToken, "", collectionID, filterJob.Dataset.DatasetID)
	}()
	go func() {
		defer wg.Done()
		releaseDate, rdErr = getReleaseDate(ctx, dc, accessToken, collectionID, filterJob.Dataset.DatasetID, filterJob.Dataset.Edition, strconv.Itoa(filterJob.Dataset.Version))
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
	if dsErr != nil {
		log.Error(ctx, "failed to get dataset", pErr, log.Data{
			"dataset_id": filterJob.Dataset.DatasetID,
		})
		setStatusCode(req, w, dsErr)
		return
	}
	if rdErr != nil {
		log.Error(ctx, "failed to get dataset release date", pErr, log.Data{
			"dataset_id": filterJob.Dataset.DatasetID,
			"edition":    filterJob.Dataset.Edition,
			"version":    strconv.Itoa(filterJob.Dataset.Version),
		})
		setStatusCode(req, w, rdErr)
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
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
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

		options = append(options, option)
	}

	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetCoverage(req, basePage, lang, filterID, geogLabel, q, pq, p, c, dimension, geogID, releaseDate, datasetDetails, areas, options, parents, hasFilterByParent, isValidationError)
	rc.BuildPage(w, m, "coverage")
}

// getAreas is a helper function that returns the GetAreasResponse or an error
func getAreas(pc PopulationClient, ctx context.Context, accessToken, popType, areaTypeID, query string) (population.GetAreasResponse, error) {
	areas, err := pc.GetAreas(ctx, population.GetAreasInput{
		AuthTokens: population.AuthTokens{
			UserAuthToken: accessToken,
		},
		PaginationParams: population.PaginationParams{
			Limit: 1000,
		},
		PopulationType: popType,
		AreaTypeID:     areaTypeID,
		Text:           url.QueryEscape(strings.TrimSpace(query)),
	})
	return areas, err
}

func getReleaseDate(ctx context.Context, dc DatasetClient, userAuthToken, collectionID, datasetID, edition, versionID string) (string, error) {
	var vErr error
	var version, initialVersion dataset.Version

	version, vErr = dc.GetVersion(ctx, userAuthToken, "", "", collectionID, datasetID, edition, versionID)
	if vErr != nil {
		return "", vErr
	}

	if version.Version != 1 {
		initialVersion, vErr = dc.GetVersion(ctx, userAuthToken, "", "", collectionID, datasetID, edition, "1")
		if vErr != nil {
			return "", vErr
		}
		return initialVersion.ReleaseDate, nil
	} else {
		return version.ReleaseDate, nil
	}
}
