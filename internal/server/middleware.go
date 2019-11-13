package server

import (
	"net/http"

	"encoding/json"
	"github.com/pkg/errors"
)

// TokenCheck asserts if request allows for asynchronous response
type TokenCheck struct{}

// ServeHTTP handling asynchronous HTTP requests in Open Service Broker Api
func (TokenCheck) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body := make(map[string]string, 0)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, errors.Wrapf(err, "while decoding body"))
		return
	}

	token := body["token"]
	if !IsValid(token) {
		writeErrorResponse(rw, http.StatusBadRequest, errors.Wrapf(err, "while verifying the token"))
		return
	}

	next(rw, r)
}
