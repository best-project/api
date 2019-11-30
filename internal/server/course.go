package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

func (srv *Server) getCourseData(body io.ReadCloser) (*internal.Course, error) {
	courseDTO := &internal.CourseDTO{}
	err := json.NewDecoder(body).Decode(&courseDTO)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding body")
	}
	return srv.converter.CourseConverter.ToModel(courseDTO), nil
}

func (srv *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	course, err := srv.getCourseData(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while decoding json body"))
		return
	}
	if err := srv.db.Course.SaveCourse(course); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while saving course"))
		return
	}

	srv.writeResponseCode(w, http.StatusCreated)
}

func (srv *Server) getAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := srv.db.Course.GetAll()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while saving course"))
		return
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}

func (srv *Server) getUserCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing jwt token"))
		return
	}
	fmt.Println(user.ID)
	courses, err := srv.db.Course.GetByUserID(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while saving course"))
		return
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}

func (srv *Server) getCoursesByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("provide user id"))
		return
	}

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing course ID: %s", id))
		return
	}

	courses, err := srv.db.Course.GetByUserID(uint(parsedID))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while saving course"))
		return
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}

func (srv *Server) getCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing course ID: %s", id))
		return
	}

	course, err := srv.db.Course.GetByID(uint(parsedID))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting from storage course ID: %s", id))
		return
	}

	dto, err := srv.converter.CourseConverter.ToDTO(course)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}
