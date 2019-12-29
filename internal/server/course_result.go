package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/server/pretty"
	"github.com/choria-io/go-validator/enum"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"sort"
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

	result := srv.converter.CourseResultConverter.ToModel(courseDTO)
	if courseDTO.Phase == internal.StartedPhase {
		results, err := srv.db.CourseResult.ListStartedForUser(courseDTO.UserID)
		if err != nil {
			srv.logger.Errorln(errors.Wrapf(err, "while listing for user %d", courseDTO.UserID))
			writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResults))
			return
		}
		if len(results) == 1 {
			result.Model = results[0].Model
		}
		if err := srv.db.CourseResult.SaveResult(result); err != nil {
			srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
			writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.CourseResult))
			return
		}
		writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.CourseResult))
		return
	}

	courseResult, bestPoints, err := srv.db.CourseResult.ReplaceIfExist(result)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while replacing started course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResult))
		return
	}

	passed, err := srv.courseLogic.CheckResult(courseResult)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while checking result"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	srv.logger.Info("trying to save result")

	dto := struct {
		Passed     bool   `json:"passed"`
		CourseId   string `json:"courseId"`
		Points     uint   `json:"points"`
		BestPoints uint   `json:"bestPoints"`
	}{}
	dto.Passed = passed
	dto.Points = courseResult.Points
	dto.CourseId = courseResult.CourseID
	dto.BestPoints = bestPoints

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

	mappedTasks := srv.db.Task.MapTasksForCourses(result)
	for _, course := range result {
		course.Task = mappedTasks[course.CourseID]
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

	mappedTasks := srv.db.Task.MapTasksForCourses(result)
	for _, course := range result {
		course.Task = mappedTasks[course.CourseID]
	}

	dto, err := srv.converter.CourseConverter.ManyToDTO(result)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) courseRanking(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		srv.logger.Errorln(errors.New("id is empty"))
		writeMessageResponse(w, http.StatusBadRequest, "provide Course id")
		return
	}
	courseResults, err := srv.db.CourseResult.ListResultsForCourse(id)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while listing courseResults"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Users))
		return
	}

	results, err := srv.converter.CourseResultConverter.ManyToDTO(courseResults)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while converting courseResults"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.CourseResults))
		return
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Points < results[j].Points
	})

	writeResponseJson(w, http.StatusOK, results)
}

func (srv *Server) usersRanking(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	users, err := srv.db.User.GetAll()
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while listing users"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Users))
		return
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Points < users[j].Points
	})

	writeResponseJson(w, http.StatusOK, srv.converter.ManyToUserStat(users))
}

func (srv *Server) userRanking(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	token, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	results, err := srv.db.CourseResult.ListBestResultsForUser(token.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while listing finished results"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResults))
		return
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Points < results[j].Points
	})

	dto, err := srv.converter.CourseResultConverter.ManyToDTO(results)
	if err != nil {
		srv.logger.Errorln(errors.New("cannot convert results to dto"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorConvert(pretty.CourseResult))
		return
	}
	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) userRankingPerCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]

	token, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	result, err := srv.db.CourseResult.GetBestResultForUserForCourse(token.ID, id)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while listing finished result"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.CourseResults))
		return
	}

	dto, err := srv.converter.CourseResultConverter.ToDTO(result)
	if err != nil {
		srv.logger.Errorln(errors.New("cannot convert result to dto"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorConvert(pretty.CourseResult))
		return
	}
	writeResponseJson(w, http.StatusOK, dto)
}
