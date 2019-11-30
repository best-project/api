package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/choria-io/go-validator/enum"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func (srv *Server) getCourseResultData(body io.ReadCloser) (*internal.CourseResult, error) {
	courseDTO := &internal.CourseResultDTO{}
	err := json.NewDecoder(body).Decode(&courseDTO)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding body")
	}
	if _, err := enum.ValidateString(courseDTO.Phase, coursePhases); err != nil {
		return nil, errors.Wrap(err, "while validate field: phase")
	}

	return srv.converter.CourseResultConverter.ToModel(courseDTO), nil
}

var coursePhases = []string{internal.StartedPhase, internal.FinishedPhase}

func (srv *Server) saveResult(w http.ResponseWriter, r *http.Request) {
	course, err := srv.getCourseResultData(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while decoding json body"))
		return
	}
	if course.Phase == internal.StartedPhase {
		if err := srv.db.CourseResult.SaveResult(course); err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while saving course"))
			return
		}
		srv.writeResponseCode(w, http.StatusCreated)
		return
	}
	course.Passed, err = srv.courseLogic.CheckResult(course)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while checking result"))
		return
	}
	srv.logger.Info("trying to save result")

	dto := struct {
		Passed bool `json:"passed"`
	}{}
	dto.Passed = course.Passed

	writeResponseObject(w, http.StatusCreated, dto)
}

func (srv *Server) getFinishedCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing jwt token"))
		return
	}
	courses, err := srv.db.CourseResult.ListFinishedForUser(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while listing courses"))
		return
	}

	result, err := srv.db.Course.GetManyByID(srv.converter.CourseResultConverter.FetchIDs(courses))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting courses"))
		return
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(result)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}

func (srv *Server) getStartedCourses(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing jwt token"))
		return
	}
	courses, err := srv.db.CourseResult.ListStartedForUser(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while listing courses"))
		return
	}

	result, err := srv.db.Course.GetManyByID(srv.converter.CourseResultConverter.FetchIDs(courses))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting courses"))
		return
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(result)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while converting courses to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, dto)
}
