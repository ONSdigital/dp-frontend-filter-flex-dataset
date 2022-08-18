package handlers

import (
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}
