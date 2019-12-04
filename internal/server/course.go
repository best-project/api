package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/pretty"
	"github.com/choria-io/go-validator/enum"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	//"github.com/olahol/go-imageupload"
	"io"
	"net/http"
	"strconv"
)

var courseTypes = []string{internal.ImageType, internal.PuzzleType}

func (srv *Server) getCourseData(body io.ReadCloser) (*internal.CourseDTO, error) {
	courseDTO := &internal.CourseDTO{}
	err := json.NewDecoder(body).Decode(&courseDTO)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding body")
	}
	if _, err := enum.ValidateString(courseDTO.Type, courseTypes); err != nil {
		return nil, errors.Wrapf(err, "invalid course type %s", courseDTO.Type)
	}

	return courseDTO, nil
}

func (srv *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	course, err := srv.getCourseData(r.Body)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while decoding json body"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewDecodeError(pretty.Course))
		return
	}
	if err := srv.validator.Struct(course); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
		return
	}
	//img, err := imageupload.Process(r, "img")
	//if err != nil {
	//	e := err.(validator.ValidationErrors)
	//	writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
	//	return
	//}

	course.MaxPoints = len(course.Data) * srv.xpForTask

	if err := srv.db.Course.SaveCourse(srv.converter.CourseConverter.ToModel(course)); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.Course))
}

func (srv *Server) getAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := srv.db.Course.GetAll()
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting all courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Courses))
		return
	}

	for _, course := range courses {
		tasks := srv.db.Task.GetTasksForCourse(course)
		course.Task = tasks
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) getUserCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, "jwt parse error")
		return
	}
	courses, err := srv.db.Course.GetByUserID(user.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.Course))
		return
	}
	for _, course := range courses {
		course.Task = srv.db.Task.GetTasksForCourse(course)
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) getCoursesByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		srv.logger.Errorln(errors.New("provide user id"))
		writeMessageResponse(w, http.StatusBadRequest, "provide user id")
		return
	}

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing course ID: %s", id))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	courses, err := srv.db.Course.GetByUserID(uint(parsedID))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.Course))
		return
	}
	for _, course := range courses {
		course.Task = srv.db.Task.GetTasksForCourse(course)
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) getCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		srv.logger.Errorln(errors.New("id is empty"))
		writeMessageResponse(w, http.StatusBadRequest, "provide Course id")
		return
	}

	course, err := srv.db.Course.GetByID(id)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting from storage course ID: %s", id))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.Course))
		return
	}
	tasks := srv.db.Task.GetTasksForCourse(course)
	course.Task = tasks

	dto, err := srv.converter.CourseConverter.ToDTO(course)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Course))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}
