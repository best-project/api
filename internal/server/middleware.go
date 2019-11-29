package server

import (
	"net/http"

	"github.com/pkg/errors"
)

// TokenCheckMiddleware asserts if request allows for asynchronous response
type TokenCheckMiddleware struct{}

// ServeHTTP handling asynchronous HTTP requests in Open Service Broker Api
func (TokenCheckMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token := r.Header.Get("Authorization")
	if token == "" || !IsValid(token) {
		writeErrorResponse(rw, http.StatusBadRequest, errors.New("while verifying auth token"))
		return
	}

	next(rw, r)
}
