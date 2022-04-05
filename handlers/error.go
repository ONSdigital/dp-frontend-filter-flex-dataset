package handlers

import (
	"net/http"
)

// validationErr is a error which occurred while validating client input.
type validationErr struct {
	error
}

func (f validationErr) Code() int {
	return http.StatusBadRequest
}
