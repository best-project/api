package server

import (
	"net/http"

	"github.com/pkg/errors"
)

// TokenCheckMiddleware asserts if request allows for asynchronous response
type TokenCheckMiddleware struct{}

// ServeHTTP handling asynchronous HTTP requests in Open Service Broker Api
func (TokenCheckMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	enableCors(&rw)
	token := r.Header.Get("Authorization")
	if _, err := ParseJWT(token); err != nil {
		writeErrorResponse(rw, http.StatusForbidden, errors.Wrap(err, "while verifying auth token"))
		return
	}

	next(rw, r)
}
