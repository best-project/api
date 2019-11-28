package server

import (
	"net/http"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/tests/service-catalog/cmd/env-tester/dto"
)

// TokenCheckMiddleware asserts if request allows for asynchronous response
type TokenCheckMiddleware struct{}

// ServeHTTP handling asynchronous HTTP requests in Open Service Broker Api
func (TokenCheckMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !IsValid(r.Header.Get("Authorization")) {
		writeErrorResponse(rw, http.StatusBadRequest, errors.New("while verifying the token"))
		return
	}

	next(rw, r)
}
