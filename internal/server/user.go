package server

import (
	"net/http"
	"fmt"
	"github.com/best-project/api/internal/model"
)

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("username")
	pass := r.URL.Query().Get("pass")

	fmt.Println(model.User{Username: name, Password: pass})

	srv.writeResponseCode(w, http.StatusOK)
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("username")
	pass := r.URL.Query().Get("pass")

	fmt.Println(model.User{Username: name, Password: pass})

	srv.writeResponseCode(w, http.StatusOK)
}
