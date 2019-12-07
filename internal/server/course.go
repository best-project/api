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
	"github.com/olahol/go-imageupload"
	"io"
	"net/http"
	"sort"
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

func (srv *Server) getTasksData(body io.ReadCloser) ([]internal.TaskDTO, error) {
	tasksDTO := make([]internal.TaskDTO, 0)
	err := json.NewDecoder(body).Decode(&tasksDTO)
	if err != nil {
		return nil, errors.Wrapf(err, "while decoding body")
	}

	return tasksDTO, nil
}

func (srv *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	courseDTO, err := srv.getCourseData(r.Body)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while decoding json body"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewDecodeError(pretty.Course))
		return
	}
	if err := srv.validator.Struct(courseDTO); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
		return
	}
	if courseDTO.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Course))
		return
	}
	img, err := imageupload.Process(r, "image")
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while fetching image"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.Course))
		return
	}
	img, err = img.ThumbnailJPEG(300, 300, 1)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting image"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.Course))
		return
	}
	courseDTO.Image = img.DataURI()
	courseDTO.MaxPoints = len(courseDTO.Data) * srv.xpForTask

	course := srv.converter.CourseConverter.ToModel(courseDTO)
	course.Task = srv.db.Task.GetTasksForCourse(course)

	if err := srv.db.Course.SaveCourse(course); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.Course))
}

func (srv *Server) updateCourse(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
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
	if !srv.db.Course.Exist(course.CourseID) {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewNotFoundError(pretty.Course))
		return
	}
	if course.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Course))
		return
	}
	img, err := imageupload.Process(r, "image")
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while fetching image"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.Course))
		return
	}
	img, err = img.ThumbnailJPEG(300, 300, 1)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting image"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.Course))
		return
	}
	course.Image = img.DataURI()
	course.MaxPoints = len(course.Data) * srv.xpForTask

	if err := srv.db.Course.SaveCourse(srv.converter.CourseConverter.ToModel(course)); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.Course))
}

func (srv *Server) addTasksToCourse(w http.ResponseWriter, r *http.Request) {
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	tasks, err := srv.getTasksData(r.Body)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while decoding json body"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewDecodeError(pretty.Course))
		return
	}
	if err := srv.validator.Struct(tasks); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
		return
	}
	if len(tasks) == 0 {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.Course))
		return
	}
	course, err := srv.db.Course.GetByID(tasks[0].CourseID)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewNotFoundError(pretty.Course))
		return
	}
	if course.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Course))
		return
	}
	course.Task = append(course.Task, srv.converter.CourseConverter.TaskConverter.ManyToModel(tasks)...)
	course.MaxPoints = len(course.Task) * srv.xpForTask

	if err := srv.db.Course.SaveCourse(course); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.Tasks))
}
func (srv *Server) removeTasksFromCourse(w http.ResponseWriter, r *http.Request) {
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

	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	task, err := srv.db.Task.GetByID(uint(parsedID))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting task"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.Task))
		return
	}
	course, err := srv.db.Course.GetByID(task.CourseID)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewNotFoundError(pretty.Course))
		return
	}
	if course.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Task))
		return
	}
	if err := srv.db.Task.RemoveByID(task); err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorRemove(pretty.Task))
		return
	}
	course.MaxPoints -= srv.xpForTask

	if err := srv.db.Course.SaveCourse(course); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewRemoveMessage(pretty.Tasks))
}

func (srv *Server) getAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := srv.db.Course.GetAll()
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting all courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Courses))
		return
	}

	mappedTasks := srv.db.Task.MapTasksForCourses(courses)
	for _, course := range courses {
		course.Task = mappedTasks[course.CourseID]
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
	mappedTasks := srv.db.Task.MapTasksForCourses(courses)
	for _, course := range courses {
		course.Task = mappedTasks[course.CourseID]
	}
	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) getUserCoursesMetadata(w http.ResponseWriter, r *http.Request) {
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
	mappedTasks := srv.db.Task.MapTasksForCourses(courses)
	for _, course := range courses {
		course.Task = mappedTasks[course.CourseID]
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
	course.Task = srv.db.Task.GetTasksForCourse(course)

	dto, err := srv.converter.CourseConverter.ToDTO(course)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Course))
		return
	}

	writeResponseJson(w, http.StatusOK, dto)
}

func (srv *Server) courseRanking(w http.ResponseWriter, r *http.Request) {
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
