package server

import (
	"encoding/json"
	"net/http"

	"github.com/svartlfheim/ymir/internal/registry"
)

type FieldError struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ErrorResponse struct {
	FieldErrors []FieldError `json:"errors"`
}

func (resp *ErrorResponse) Add(name string, value interface{}) {
	resp.FieldErrors = append(resp.FieldErrors, FieldError{
		Name:  name,
		Value: value,
	})
}

func handleValidationErrorsResponse(errs []registry.ValidationError, code int, w http.ResponseWriter) {
	errResp := ErrorResponse{}

	for _, err := range errs {
		errResp.Add(err.Field, err.Message)
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errResp)
}

type Meta struct {
	Ref string `json:"ref"`
}

type ResourceResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func handleResourceResponse(r interface{}, code int, w http.ResponseWriter) {
	resp := ResourceResponse{
		Meta: Meta{
			Ref: "something",
		},
		Data: r,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
