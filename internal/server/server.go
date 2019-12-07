package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal/converter"
	"github.com/best-project/api/internal/service"
	"github.com/best-project/api/internal/storage"
	"github.com/go-playground/validator"
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

	courseLogic *service.CourseLogic
	validator   *validator.Validate

	host string

	xpForTask int

	converter *converter.Converter

	fbCallbackURL        string
	instagramCallbackURL string
}

const (
	fbEndpoint        = "/register/fb/callback"
	instagramEndpoint = "/register/instagram/callback"
)

func NewServer(db *storage.Database, fb facebook.Interface, courseLogic *service.CourseLogic, logger *logrus.Logger) *Server {
	host := "https://secret-cliffs-62699.herokuapp.com/"

	return &Server{
		logger: logger,
		fb:     fb,
		db:     db,

		converter:   converter.NewConverter(db),
		courseLogic: courseLogic,
		validator:   validator.New(),

		xpForTask: 10,

		host:                 host,
		fbCallbackURL:        fmt.Sprintf("%s%s", host, fbEndpoint),
		instagramCallbackURL: fmt.Sprintf("%s%s", host, instagramEndpoint),
	}
}

// Handle creates an http handler
func (srv *Server) Handle() http.Handler {
	var rtr = mux.NewRouter()

	tokenCheckMiddleware := TokenCheckMiddleware{}

	rtr.Path("/login").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.loginUser)))
	rtr.Path("/register/user").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.createUser)))
	//rtr.Path("/register/fb").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.redirectToFb)))
	//rtr.Path("/register/instagram").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.redirectToInstagram)))

	rtr.Path("/user").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserByToken)))
	rtr.Path("/user/{name}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUser)))
	rtr.Path("/user/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserByID)))
	rtr.Path("/user/courses/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getCoursesByUserID)))

	rtr.Path("/course/create").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.createCourse)))
	rtr.Path("/course/update").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.updateCourse)))
	rtr.Path("/course/task/add").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.addTasksToCourse)))
	rtr.Path("/course/task/remove").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.removeTasksFromCourse)))

	rtr.Path("/course/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getCourse)))
	rtr.Path("/courses").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserCourses)))
	rtr.Path("/courses/meta").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserCoursesMetadata)))
	rtr.Path("/courses/all").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getAllCourses)))
	rtr.Path("/courses/started").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getStartedCourses)))
	rtr.Path("/courses/finished").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getFinishedCourses)))

	rtr.Path("/result/save").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.saveResult)))
	rtr.Path("/ranking").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.fetchByXP)))
	rtr.Path("/ranking/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.courseRanking)))

	rtr.Path("/status").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.statusHandler)))

	return rtr
}

func (srv *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(http.StatusOK)))
}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	dto := struct {
		Message string
		Code    int
	}{
		Message: err.Error(),
		Code:    code,
	}
	writeResponseJson(w, code, dto)
	return
}

func writeMessageResponse(w http.ResponseWriter, code int, msg string) {
	dto := struct {
		Message string
		Code    int
	}{
		Message: msg,
		Code:    code,
	}
	writeResponseJson(w, code, dto)
	return
}

func writeResponseJson(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
