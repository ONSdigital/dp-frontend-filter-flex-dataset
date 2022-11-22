package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/mapper"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/gorilla/mux"
)

// ChangeDimensions Handler
func ChangeDimensions(rc RenderClient) http.HandlerFunc {
	return handlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		changeDimensions(w, req, rc, accessToken, collectionID, lang)
	})
}

func changeDimensions(w http.ResponseWriter, req *http.Request, rc RenderClient, accessToken, collectionID, lang string) {
	vars := mux.Vars(req)
	fid := vars["filterID"]
	basePage := rc.NewBasePageModel()
	m := mapper.CreateGetChangeDimensions(req, basePage, lang, fid)
	rc.BuildPage(w, m, "dimensions")
}
