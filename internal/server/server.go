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
	"regexp"
	"strconv"
)

type Server struct {
	logger logrus.FieldLogger
	fb     facebook.Interface
	db     *storage.Database

	courseLogic service.CourseResultHandler
	validator   *validator.Validate

	host      string
	xpForTask int

	converter *converter.Converter
}

func NewServer(db *storage.Database, fb facebook.Interface, courseLogic *service.CourseLogic, logger *logrus.Logger) *Server {
	return &Server{
		logger: logger,
		fb:     fb,
		db:     db,

		converter:   converter.NewConverter(),
		courseLogic: courseLogic,
		validator:   validator.New(),

		xpForTask: 10,
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Handle creates an http handler
func (srv *Server) Handle() http.Handler {
	var rtr = mux.NewRouter()
	tokenCheckMiddleware := TokenCheckMiddleware{}

	rtr.Path("/login").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.loginUser)))
	rtr.Path("/register/user").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.createUser)))

	rtr.Path("/user").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserByToken)))
	rtr.Path("/user/update").Methods(http.MethodPut).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.updateUser)))
	rtr.Path("/user/refresh/{token}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.refreshToken)))
	rtr.Path("/user/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserByID)))
	rtr.Path("/user/courses/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getCoursesByUserID)))
	rtr.Path("/users/ranking").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.userRanking)))
	rtr.Path("/users/ranking/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.userRankingPerCourse)))

	rtr.Path("/course/create").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.createCourse)))
	rtr.Path("/course/update").Methods(http.MethodPut).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.updateCourse)))
	rtr.Path("/course/task/add").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.addTasksToCourse)))
	rtr.Path("/course/task/edit").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.editTask)))
	rtr.Path("/course/task/remove/{id}").Methods(http.MethodDelete).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.removeTasksFromCourse)))

	rtr.Path("/course/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getCourse)))
	rtr.Path("/courses").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserCourses)))
	rtr.Path("/courses/meta").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getUserCoursesMetadata)))
	rtr.Path("/courses/all").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getAllCourses)))
	rtr.Path("/courses/all/meta").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getAllCoursesMetadata)))
	rtr.Path("/courses/started").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getStartedCourses)))
	rtr.Path("/courses/finished").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.getFinishedCourses)))

	rtr.Path("/result/save").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.saveResult)))
	rtr.Path("/rate").Methods(http.MethodPost).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.rateCourse)))
	rtr.Path("/ranking").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.usersRanking)))
	rtr.Path("/ranking/{id}").Methods(http.MethodGet).Handler(negroni.New(tokenCheckMiddleware, negroni.WrapFunc(srv.courseRanking)))

	rtr.Path("/images/{name}").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.getImage)))
	rtr.Path("/status").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.statusHandler)))

	return rtr
}

func (srv *Server) getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["name"]

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

	http.ServeFile(w, r, fmt.Sprintf("images/%s", filename))
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

func (srv *Server) ParseFormCollection(r *http.Request, typeName string) []map[string]string {
	var result []map[string]string
	r.ParseForm()
	for key, values := range r.Form {
		re := regexp.MustCompile(typeName + "\\[([0-9]+)\\]\\[([a-zA-Z]+)\\]")
		matches := re.FindStringSubmatch(key)

		if len(matches) >= 3 {

			index, _ := strconv.Atoi(matches[1])

			for index >= len(result) {
				result = append(result, map[string]string{})
			}

			result[index][matches[2]] = values[0]
		}
	}
	return result
}
