package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal/converter"
	"github.com/best-project/api/internal/storage"
	"github.com/gorilla/mux"
	"github.com/madebyais/facebook-go-sdk"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

type Server struct {
	logger logrus.FieldLogger
	fb     facebook.Interface
	db     *storage.Database

	host string

	converter *converter.Converter

	fbCallbackURL        string
	instagramCallbackURL string
}

const (
	fbEndpoint        = "/register/fb/callback"
	instagramEndpoint = "/register/instagram/callback"
)

func NewServer(db *storage.Database, fb facebook.Interface) *Server {
	host := "https://secret-cliffs-62699.herokuapp.com/"

	return &Server{
		logger: logrus.New(),
		fb:     fb,
		db:     db,

		converter: converter.NewConverter(),

		host:                 host,
		fbCallbackURL:        fmt.Sprintf("%s%s", host, fbEndpoint),
		instagramCallbackURL: fmt.Sprintf("%s%s", host, instagramEndpoint),
	}
}

// Handle creates an http handler
func (srv *Server) Handle() http.Handler {
	var rtr = mux.NewRouter()

	rtr.Path("/register/user").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.createUser)))

	rtr.Path("/register/fb").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.redirectToFb)))
	rtr.Path(fbEndpoint).Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.createUserFb)))

	rtr.Path("/register/instagram").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.redirectToInstagram)))
	rtr.Path(instagramEndpoint).Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.createUserInstagram)))

	rtr.Path("/course/create").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.createCourse)))

	rtr.Path("/user/{name}").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.getUser)))
	rtr.Path("/course/{id}").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.getCourse)))

	rtr.Path("/status").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.statusHandler)))

	return rtr
}

func (srv *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	srv.writeResponseCode(w, http.StatusOK)
	return
}

func (srv *Server) writeResponseCode(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(strconv.Itoa(code)))
}

func (srv *Server) writeResponseBody(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	dto := struct {
		Message string
		Code    int
	}{
		Message: err.Error(),
		Code:    code,
	}
	writeResponseObject(w, code, dto)
	return
}

func writeResponseObject(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
