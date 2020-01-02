package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/server/pretty"
	"github.com/choria-io/go-validator/enum"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var courseTypes = []string{internal.NormalType, internal.PuzzleType}

const (
	MB = 1 << 20
)

func (srv *Server) getCourseData(r *http.Request) (*internal.CourseDTO, error) {
	courseDTO := &internal.CourseDTO{}
	if err := r.ParseMultipartForm(MB * 20); err != nil {
		return nil, err
	}
	courseDTO.Name = r.FormValue("name")

	imgPath, _ := srv.saveImage(r)
	courseDTO.Image = imgPath
	courseDTO.CourseID = r.FormValue("courseId")
	courseDTO.Type = r.FormValue("type")
	courseDTO.Language = r.FormValue("language")
	courseDTO.Description = r.FormValue("description")
	courseDTO.DifficultyLevel = r.FormValue("difficultyLevel")

	if _, err := enum.ValidateString(courseDTO.Type, courseTypes); err != nil {
		return nil, errors.Wrapf(err, "invalid course type %s", courseDTO.Type)
	}
	data := srv.ParseFormCollection(r, "data")
	tasksDTO := make([]internal.TaskDTO, 0)
	taskDTO := internal.TaskDTO{}

	for _, task := range data {
		taskDTO.Image = task["image"]
		taskDTO.Word = task["word"]
		taskDTO.Translate = task["translate"]
		tasksDTO = append(tasksDTO, taskDTO)
	}
	courseDTO.Data = tasksDTO

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
	enableCors(&w)
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	course, err := srv.getCourseData(r)
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

	course.Rate = 3
	course.UserID = user.ID
	course.MaxPoints = len(course.Data) * srv.xpForTask

	if err := srv.db.Course.SaveCourse(srv.converter.CourseConverter.ToModel(course), srv.xpForTask); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.Course))
		return
	}

	writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.Course))
}

func (srv *Server) updateCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	course, err := srv.getCourseData(r)
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

	courseModel, err := srv.db.Course.GetByID(course.CourseID)
	if err != nil {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewErrorGet(pretty.User))
		return
	}
	if courseModel.Description != "" {
		courseModel.Description = course.Description
	}
	if courseModel.Image != "" {
		courseModel.Image = course.Image
	}
	if courseModel.Language != "" {
		courseModel.Language = course.Language
	}
	if courseModel.Name != "" {
		courseModel.Name = course.Name
	}
	courseModel.UpdatedAt = time.Now()

	if err := srv.db.Course.SaveCourse(courseModel, srv.xpForTask); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.Course))
}

func (srv *Server) rateCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	rateDTO := struct {
		Rate     int    `json:"rate"`
		CourseID string `json:"courseId"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&rateDTO)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	course, err := srv.db.Course.GetByID(rateDTO.CourseID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting course %s", rateDTO.CourseID))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.Course))
		return
	}

	rate := float32(course.RateCounter) * course.Rate
	rate += float32(rateDTO.Rate)
	course.Rate = rate / float32(course.RateCounter+1)
	course.RateCounter++

	if err := srv.db.Course.SaveCourse(course, srv.xpForTask); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewCreateMessage(pretty.Course))
}

func (srv *Server) saveImage(r *http.Request) (string, error) {
	//Content-Type
	file, _, err := r.FormFile("image")
	if err != nil {
		return "", errors.Wrap(err, "while getting image from form")
	}
	defer file.Close()

	imgPath := fmt.Sprintf("/images/%s.png", xid.New().String())
	saveFile, err := os.Create(imgPath)
	if err != nil {
		return "", err
	}
	defer saveFile.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(saveFile, file)
	if err != nil {
		return "", err
	}
	return imgPath, nil
}

func (srv *Server) addTasksToCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	if err := r.ParseMultipartForm(MB * 50); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing multipart file"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	taskDTO := internal.TaskDTO{}
	imgPath, err := srv.saveImage(r)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing multipart file"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.Task))
		return
	}
	taskDTO.Image = imgPath
	taskDTO.CourseID = r.FormValue("courseId")
	taskDTO.Word = r.FormValue("word")
	taskDTO.Translate = r.FormValue("translate")
	if err := srv.validator.Struct(taskDTO); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
		return
	}

	course, err := srv.db.Course.GetByID(taskDTO.CourseID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting by id %s", taskDTO.CourseID))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewNotFoundError(pretty.Course))
		return
	}
	if course.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Course))
		return
	}
	course.Task = append(course.Task, srv.converter.CourseConverter.TaskConverter.ConvertToModel(taskDTO))
	course.MaxPoints = len(course.Task) * srv.xpForTask

	if err := srv.db.Course.SaveCourse(course, srv.xpForTask); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.Course))
		return
	}
	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.Tasks))
}

func (srv *Server) editTask(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	user, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	if err := r.ParseMultipartForm(MB * 50); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing multipart file"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	id := r.FormValue("id")
	if id == "" {
		srv.logger.Info("no task ID provided")
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	taskDTO := internal.TaskDTO{}
	imgPath, err := srv.saveImage(r)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing multipart file"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.Task))
		return
	}
	taskDTO.ID = id
	taskDTO.Image = imgPath
	taskDTO.CourseID = r.FormValue("courseId")
	taskDTO.Word = r.FormValue("word")
	taskDTO.Translate = r.FormValue("translate")
	if err := srv.validator.Struct(taskDTO); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.Course, e))
		return
	}
	course, err := srv.db.Course.GetByID(taskDTO.CourseID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting by id %s", taskDTO.CourseID))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewNotFoundError(pretty.Course))
		return
	}
	if course.UserID != user.ID {
		writeMessageResponse(w, http.StatusForbidden, pretty.NewForbiddenError(pretty.Course))
		return
	}
	task := srv.converter.CourseConverter.TaskConverter.ConvertToModel(taskDTO)
	if err := srv.db.Task.SaveTask(&task); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving task"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.Task))
		return
	}
	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.Tasks))
}
func (srv *Server) removeTasksFromCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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
	course, err := srv.db.Course.GetByID(strconv.Itoa(int(task.CourseID)))
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

	if err := srv.db.Course.SaveCourse(course, srv.xpForTask); err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while saving course"))
		writeMessageResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeMessageResponse(w, http.StatusOK, pretty.NewRemoveMessage(pretty.Tasks))
}

func (srv *Server) applyBestPointsToCourses(w http.ResponseWriter, r *http.Request, courses []internal.CourseDTO) {
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
	for _, course := range courses {
		for _, result := range results {
			r, _ := strconv.Atoi(course.CourseID)
			if uint(r) == result.CourseID {
				course.BestPoints = int(result.Points)
			}
		}
	}

	writeResponseJson(w, http.StatusOK, courses)
}

func (srv *Server) getAllCourses(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	courses, err := srv.db.Course.GetAll()
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting all courses"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorList(pretty.Courses))
		return
	}
	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	srv.applyBestPointsToCourses(w, r, dto)
}

func (srv *Server) getUserCourses(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	srv.applyBestPointsToCourses(w, r, dto)
}

func (srv *Server) getAllCoursesMetadata(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	courses, err := srv.db.Course.GetAll()
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

	srv.applyBestPointsToCourses(w, r, dto)
}

func (srv *Server) getUserCoursesMetadata(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	srv.applyBestPointsToCourses(w, r, dto)
}

func (srv *Server) getCoursesByUserID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	dto, err := srv.converter.CourseConverter.ManyToDTO(courses)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while converting courses to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.Courses))
		return
	}

	srv.applyBestPointsToCourses(w, r, dto)
}

func (srv *Server) getCourse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	srv.applyBestPointsToCourses(w, r, []internal.CourseDTO{*dto})
}
