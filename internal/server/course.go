package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (srv *Server) getCourseData(w http.ResponseWriter, r *http.Request) *internal.Course {
	courseDTO := &internal.CourseDTO{}
	err := json.NewDecoder(r.Body).Decode(&courseDTO)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while decoding body"))
		return nil
	}
	return srv.converter.CourseConverter.ToModel(courseDTO)
}

func (srv *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	srv.db.Course.SaveCourse(srv.getCourseData(w, r))

	srv.writeResponseCode(w, http.StatusCreated)
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

	writeResponseObject(w, http.StatusOK, course)
}
