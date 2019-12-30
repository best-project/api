package server

import (
	"encoding/json"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/server/pretty"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (srv *Server) readUserData(body io.ReadCloser) (*internal.UserDTO, error) {
	userDTO := &internal.UserDTO{}
	err := json.NewDecoder(body).Decode(&userDTO)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding user")
	}

	return userDTO, nil
}

func (srv *Server) getUserByToken(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	token, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	user, err := srv.db.User.GetByID(token.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting user with name: %s", token.Email))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.User))
		return
	}

	result, err := srv.converter.ToDTO(*user)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while converting to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.User))
		return
	}

	writeResponseJson(w, http.StatusOK, result)
}

func (srv *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing course ID: %s", id))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	user, err := srv.db.User.GetByID(uint(parsedID))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting user with id: %s", id))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.User))
		return
	}

	result, err := srv.converter.ToDTO(*user)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while converting to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.User))
		return
	}

	writeResponseJson(w, http.StatusOK, result)
}

func (srv *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userDTO := &struct {
		Email    string `json:"email" validate:"required,email,max=250"`
		Password string `json:"password" validate:"required,min=8,max=250"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&userDTO)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.User))
		return
	}
	if err := srv.validator.Struct(userDTO); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.User, e))
		return
	}

	srv.logger.Infof("trying to login user %s", userDTO.Email)
	users, err := srv.db.User.GetByMail(userDTO.Email)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting by mail %s", userDTO.Email))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	if len(users) == 0 {
		srv.logger.Errorln(errors.New("user not found"))
		writeMessageResponse(w, http.StatusNotFound, pretty.NewNotFoundError(pretty.User))
		return
	}
	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while comparing passwords"))
		writeMessageResponse(w, http.StatusForbidden, pretty.NewInternalError())
		return
	}

	now := time.Now()
	token, err := NewJWT(NewCustomPayload(&UserClaim{ID: user.ID, Email: user.Email}, now.Add(time.Minute*30).Unix()))
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while creating token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	refreshToken, err := NewJWT(NewCustomPayload(&UserClaim{ID: user.ID, Email: user.Email}, now.Add(time.Hour*12).Unix()))
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while creating refresh token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	result, err := srv.converter.ToDTO(user)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while converting to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.User))
		return
	}
	result.Token = token
	result.RefreshToken = refreshToken

	writeResponseJson(w, http.StatusOK, result)
}

func (srv *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	token := vars["token"]
	if token == "" {
		srv.logger.Errorln("token cannot be empty")
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewForbiddenError(pretty.User))
		return
	}

	userClaim, err := ParseJWT(token)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while parsing jwt"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewForbiddenError(pretty.User))
		return
	}
	user, err := srv.db.User.GetByID(userClaim.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while parsing jwt"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.User))
		return
	}
	now := time.Now()
	token, err = NewJWT(NewCustomPayload(&UserClaim{ID: user.ID, Email: user.Email}, now.Add(time.Minute*30).Unix()))
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while creating token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	result, err := srv.converter.ToDTO(*user)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while converting to dto"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorConvert(pretty.User))
		return
	}
	result.Token = token

	writeResponseJson(w, http.StatusOK, result)
}

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userData, err := srv.readUserData(r.Body)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.User))
		return
	}
	if err := srv.validator.Struct(userData); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.User, e))
		return
	}

	user := srv.converter.ToModel(userData)
	if srv.db.User.Exist(user) {
		srv.logger.Errorln(errors.New("user already exists"))
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewAlreadyExistError(pretty.User))
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while hashing password"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}

	user.Password = string(pass)
	user.Level = 1
	user.NextLevel = srv.courseLogic.PointForNextLevel(user.Points)
	if err := srv.db.User.SaveUser(user); err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while saving user"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.User))
		return
	}

	srv.logger.Infof("user %s was created", user.Email)
	writeMessageResponse(w, http.StatusCreated, pretty.NewCreateMessage(pretty.User))
}

func (srv *Server) getUserDataFromForm(r *http.Request) (*internal.UserDTO, error) {
	imgPath, _ := srv.saveImage(r)

	userDTO := &internal.UserDTO{
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Password:  r.FormValue("password"),
		Avatar:    imgPath,
	}
	return userDTO, nil
}

func (srv *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userData, err := srv.getUserDataFromForm(r)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewDecodeError(pretty.User))
		return
	}
	if err := srv.validator.Struct(userData); err != nil {
		e := err.(validator.ValidationErrors)
		writeMessageResponse(w, http.StatusBadRequest, pretty.NewErrorValidate(pretty.User, e))
		return
	}
	token, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while parsing jwt token"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
		return
	}
	user, err := srv.db.User.GetByID(token.ID)
	if err != nil {
		srv.logger.Errorln(errors.Wrapf(err, "while getting user"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorGet(pretty.User))
		return
	}

	var pass []byte
	pass = []byte(user.Password)
	if userData.Password != "" {
		pass, err = bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			srv.logger.Errorln(errors.Wrap(err, "while hashing password"))
			writeMessageResponse(w, http.StatusInternalServerError, pretty.NewInternalError())
			return
		}
	}
	user.Password = string(pass)

	if user.Email != "" {
		user.Email = userData.Email
	}
	if user.FirstName != "" {
		user.FirstName = userData.FirstName
	}
	if user.LastName != "" {
		user.LastName = userData.LastName
	}
	if user.Avatar != "" {
		user.Avatar = userData.Avatar
	}

	if err := srv.db.User.SaveUser(user); err != nil {
		srv.logger.Errorln(errors.Wrap(err, "while saving user"))
		writeMessageResponse(w, http.StatusInternalServerError, pretty.NewErrorSave(pretty.User))
		return
	}
	writeMessageResponse(w, http.StatusOK, pretty.NewUpdateMessage(pretty.User))
}
