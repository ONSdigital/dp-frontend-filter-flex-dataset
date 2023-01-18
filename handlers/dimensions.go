package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/population"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// DimensionSelector Handler
func (f *FilterFlex) DimensionSelector() http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		dimensionSelector(w, req, f, collectionID, accessToken, lang)
	})
}

func dimensionSelector(w http.ResponseWriter, req *http.Request, f *FilterFlex, collectionID, accessToken, lang string) {
	ctx := req.Context()
	vars := mux.Vars(req)
	filterID := vars["filterID"]
	dimensionName := vars["name"]
	isValidationError, _ := strconv.ParseBool(req.URL.Query().Get("error"))

	logData := log.Data{
		"filter_id": filterID,
	}

	eb, serviceMsg, err := getZebContent(ctx, f.ZebedeeClient, accessToken, collectionID, lang)
	// log zebedee error but don't set a server error
	if err != nil {
		log.Error(ctx, "unable to get homepage content", err, log.Data{"homepage_content": err})
	}

	currentFilter, _, err := f.FilterClient.GetJobState(ctx, accessToken, "", "", collectionID, filterID)
	if err != nil {
		log.Error(ctx, "failed to get job state", err, logData)
		setStatusCode(req, w, err)
		return
	}

	filterDimension, _, err := f.FilterClient.GetDimension(ctx, accessToken, "", collectionID, filterID, dimensionName)
	if err != nil {
		log.Error(ctx, "failed to find dimension in filter", err, logData)
		setStatusCode(req, w, err)
		return
	}

	basePage := f.Render.NewBasePageModel()

	if !isAreaType(filterDimension) {
		isMultivariate, err := isMultivariateDataset(ctx, f.DatasetClient, accessToken, collectionID, currentFilter.Dataset.DatasetID)
		if err != nil {
			log.Error(ctx, "failed to determine if filter is multivariate", err, log.Data{
				"filter_id":  filterID,
				"dataset_id": currentFilter.Dataset.DatasetID,
			})
			setStatusCode(req, w, err)
			return
		}
		if !isMultivariate || !f.EnableMultivariate {
			err = &clientErr{errors.New("invalid request")}
			setStatusCode(req, w, err)
			return
		}
		cats, err := f.PopulationClient.GetCategorisations(ctx, population.GetCategorisationsInput{
			AuthTokens: population.AuthTokens{
				UserAuthToken: accessToken,
			},
			PaginationParams: population.PaginationParams{
				Limit: 1000,
			},
			PopulationType: currentFilter.PopulationType,
			Dimension:      dimensionName,
		})
		if err != nil {
			log.Error(ctx, "failed to get categorisations", err, log.Data{
				"filter_id":       filterID,
				"population_type": currentFilter.PopulationType,
				"dimension":       dimensionName,
			})
			setStatusCode(req, w, err)
			return
		}

		selector := mapper.CreateCategorisationsSelector(req, basePage, filterDimension.Label, lang, filterID, dimensionName, serviceMsg, eb, cats, isValidationError)
		f.Render.BuildPage(w, selector, "selector")
		return
	}

	// The total_count is the only field required
	opts, _, err := f.FilterClient.GetDimensionOptions(ctx, accessToken, "", collectionID, filterID, dimensionName, &filter.QueryParams{Limit: 0})
	if err != nil {
		log.Error(ctx, "failed to get options for dimension", err, logData)
		setStatusCode(req, w, err)
		return
	}

	hasOpts := opts.TotalCount > 0

	areaTypes, err := f.PopulationClient.GetAreaTypes(ctx, population.GetAreaTypesInput{
		AuthTokens: population.AuthTokens{
			UserAuthToken: accessToken,
		},
		PaginationParams: population.PaginationParams{
			Limit: 1000,
		},
		PopulationType: currentFilter.PopulationType,
	})
	if err != nil {
		log.Error(ctx, "failed to get population area types", err, logData)
		setStatusCode(req, w, err)
		return
	}

	details, err := f.DatasetClient.GetVersion(ctx, accessToken, "", "", collectionID, currentFilter.Dataset.DatasetID, currentFilter.Dataset.Edition, strconv.Itoa(currentFilter.Dataset.Version))
	if err != nil {
		log.Error(ctx, "failed to get dataset version", err, log.Data{
			"dataset": currentFilter.Dataset.DatasetID,
			"edition": currentFilter.Dataset.Edition,
			"version": currentFilter.Dataset.Version,
		})
		setStatusCode(req, w, err)
		return
	}

	dataset, err := f.DatasetClient.Get(ctx, accessToken, "", collectionID, currentFilter.Dataset.DatasetID)
	if err != nil {
		log.Error(ctx, "failed to get dataset", err, log.Data{
			"dataset": currentFilter.Dataset.DatasetID,
		})
		setStatusCode(req, w, err)
		return
	}

	releaseDate, err := getReleaseDate(ctx, f.DatasetClient, accessToken, collectionID, currentFilter.Dataset.DatasetID, currentFilter.Dataset.Edition, strconv.Itoa(currentFilter.Dataset.Version))
	if err != nil {
		log.Error(ctx, "failed to get release date", err, log.Data{
			"dataset": currentFilter.Dataset.DatasetID,
			"edition": currentFilter.Dataset.Edition,
		})
		setStatusCode(req, w, err)
		return
	}

	selector := mapper.CreateAreaTypeSelector(req, basePage, lang, filterID, areaTypes.AreaTypes, filterDimension, details.LowestGeography, releaseDate, dataset, isValidationError, hasOpts, serviceMsg, eb)
	f.Render.BuildPage(w, selector, "selector")
}

// isAreaType determines if the current dimension is an area type
func isAreaType(dimension filter.Dimension) bool {
	if dimension.IsAreaType == nil {
		return false
	}

	return *dimension.IsAreaType
}
