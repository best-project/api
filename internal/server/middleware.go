package server

import (
	"net/http"

	"encoding/json"
	"github.com/pkg/errors"
)

// TokenCheckMiddleware asserts if request allows for asynchronous response
type TokenCheckMiddleware struct{}

// ServeHTTP handling asynchronous HTTP requests in Open Service Broker Api
func (TokenCheckMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	dto := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, errors.Wrapf(err, "while decoding body"))
		return
	}
	if !IsValid(dto.Token) {
		writeErrorResponse(rw, http.StatusBadRequest, errors.Wrapf(err, "while verifying the token"))
		return
	}

	next(rw, r)
}
