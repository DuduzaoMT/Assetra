package restutil

import (
	"assetra/security"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

var (
	ErrEmptyBody     = errors.New("body can't be empty")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("access forbidden")
	ErrAccountLocked = errors.New("account locked due to multiple failed login attempts")
)

type JError struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	e := "error"
	if err != nil {
		e = err.Error()
	}
	WriteAsJson(w, statusCode, JError{e})
}

func WriteAsJson(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func AuthRequestWithId(r *http.Request) (*security.TokenPayload, error) {
	payload, err := AuthRequestToken(r)
	if err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	// allow the user itself or an "admin"
	if payload.UserId != vars["id"] && !ContainsRole(payload.Roles, "admin") {
		return nil, ErrUnauthorized
	}
	return payload, nil
}

func ContainsRole(slice []string, role string) bool {
	return slices.Contains(slice, role)
}

func AuthRequestToken(r *http.Request) (*security.TokenPayload, error) {
	token, err := security.ExtractToken(r)
	if err != nil {
		return nil, err
	}
	payload, err := security.NewTokenPayload(token)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
