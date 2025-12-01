package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New(validator.WithRequiredStructEnabled())
	})

	return validatorInstance
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_568
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}
