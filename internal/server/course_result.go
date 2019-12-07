package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/pretty"
	"github.com/choria-io/go-validator/enum"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func (srv *Server) getCourseResultData(body io.ReadCloser) (*internal.CourseResultDTO, error) {
	courseDTO := &internal.CourseResultDTO{}
	err := json.NewDecoder(body).Decode(&courseDTO)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding body")
	}
	if _, err := enum.ValidateString(courseDTO.Phase, coursePhases); err != nil {
		return nil, errors.Wrap(err, "while validate field: phase")
	}

	return courseDTO, nil
}

var coursePhases = []string{internal.StartedPhase, internal.FinishedPhase}

func (srv *Server) saveResult(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	courseDTO, err := srv.getCourseResultData(r.Body)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while decoding json body"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewDecodeError(pretty.CourseResult))
		return
	}
	if err := srv.validator.Struct(courseDTO); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorValidate(pretty.CourseResult, e))
		return
	}
	if !srv.db.Course.Exist(courseDTO.CourseID) {
		srv.logger.Errorln(errors.New("course not exist"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewNotFoundError(pretty.Course))
		return
	}
	courseDTO.UserID = user.ID

	course, err := srv.db.CourseResult.ReplaceIfExist(srv.converter.CourseResultConverter.ToModel(courseDTO))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while replacing started course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResult))
		return
	}
	if course.Phase == internal.StartedPhase {
		if err := srv.db.CourseResult.SaveResult(course); err != nil {
			srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
			writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.CourseResult))
			return
		}
		writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.CourseResult))
		return
	}
	passed, err := srv.courseLogic.CheckResult(course)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while checking result"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	srv.logger.Info("trying to save result")

	dto := struct {
		Passed bool `json:"passed"`
	}{}
	dto.Passed = passed

	writeResponseJson(w, http.StatusCreated, dto)
}

func (srv *Server) getFinishedCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	courses, err := srv.db.CourseResult.ListFinishedForUser(user.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while listing courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResults))
		return
	}

	result, err := srv.db.Course.GetManyByID(srv.converter.CourseResultConverter.FetchIDs(courses))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Courses))
		return
	}

	mappedCourses := srv.db.Task.MapTasksForCourses(result)
	for _, course := range result {
		course.Task = mappedCourses[course.CourseID]
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(result)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) getStartedCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	courses, err := srv.db.CourseResult.ListStartedForUser(user.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while listing courses results"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResults))
		return
	}
	result, err := srv.db.Course.GetManyByID(srv.converter.CourseResultConverter.FetchIDs(courses))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Courses))
		return
	}

	mappedCourses := srv.db.Task.MapTasksForCourses(result)
	for _, course := range result {
		course.Task = mappedCourses[course.CourseID]
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(result)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}