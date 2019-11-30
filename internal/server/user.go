package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strconv"
)

func (srv *Server) readUserData(body io.ReadCloser) (*internal.User, error) {
	userDTO := &internal.UserDTO{}
	err := json.NewDecoder(body).Decode(&userDTO)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding user")
	}
	return srv.converter.ToModel(userDTO), nil
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	users, err := srv.db.User.GetByName(name)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting user with name: %s", name))
		return
	}
	if len(users) == 0 {
		writeErrorResponse(w, http.StatusNotFound, errors.New("user not found"))
		return
	}

	user := users[0]

	result, err := srv.converter.ToDTO(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while converting to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, result)
}

func (srv *Server) getUserByToken(w http.ResponseWriter, r *http.Request) {
	token, err := ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing jwt token"))
		return
	}

	user, err := srv.db.User.GetByID(token.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting user with name: %s", token.Email))
		return
	}

	result, err := srv.converter.ToDTO(*user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while converting to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, result)
}

func (srv *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while parsing course ID: %s", id))
		return
	}

	user, err := srv.db.User.GetByID(uint(parsedID))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting user with id: %s", id))
		return
	}

	result, err := srv.converter.ToDTO(*user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while converting to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, result)
}

func (srv *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	userData, err := srv.readUserData(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while reading user %s data", userData.Email))
		return
	}

	srv.logger.Infof("trying to login user %s", userData.Email)
	users, err := srv.db.User.GetByMail(userData.Email)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting by mail %s", userData.Email))
		return
	}
	if len(users) == 0 {
		writeErrorResponse(w, http.StatusNotFound, errors.New("user not found"))
		return
	}
	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password)); err != nil {
		writeErrorResponse(w, http.StatusForbidden, errors.Wrap(err, "while comparing passwords"))
		return
	}

	token, err := NewJWT(NewCustomPayload(&UserClaim{ID: user.ID, Email: user.Email}))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while creating token"))
		return
	}
	user.Token = token
	if err := srv.db.User.SaveUser(&user); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while saving user"))
		return
	}
	result, err := srv.converter.ToDTO(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while converting to dto"))
		return
	}

	writeResponseObject(w, http.StatusOK, result)
}

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := srv.readUserData(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while reading user %s data", user.Username))
		return
	}

	if srv.db.User.Exist(user) {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("user already exists"))
		return
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while hashing password"))
		return
	}

	user.Password = string(pass)
	if err := srv.db.User.SaveUser(user); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while saving user"))
		return
	}

	srv.logger.Infof("user %s was created", user.Email)
	srv.writeResponseCode(w, http.StatusCreated)
}

func (srv *Server) redirectToFb(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://www.facebook.com/v2.10/dialog/oauth?client_id=%s&redirect_uri=%s", srv.fb.GetAppID(), srv.fbCallbackURL), http.StatusTemporaryRedirect)
}

func (srv *Server) createUserFb(w http.ResponseWriter, r *http.Request) {

	data, err := srv.fb.GenerateAccessToken(srv.fbCallbackURL, r.URL.Query().Get("code"))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.New("while generating access token"))
		return
	}
	srv.fb.SetAccessToken(fmt.Sprintf("%v", data["access_token"]))

	feed, err := srv.fb.API("/me").Get()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while getting user info"))
		return
	}
	fmt.Println("FB:", feed)

	token, err := NewJWT(NewCustomPayload(&UserClaim{}))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.New("while creating token"))
		return
	}

	srv.writeResponseBody(w, http.StatusCreated, []byte(token))
}

func (srv *Server) redirectToInstagram(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "", http.StatusTemporaryRedirect)
}

func (srv *Server) createUserInstagram(w http.ResponseWriter, r *http.Request) {
	//user := &internal.User{}
	//
	//
	//
	//token, err := jwt.Sign(pl, signingKey)
	//if err != nil {
	//	writeErrorResponse(w, http.StatusInternalServerError, "while creating user")
	//	return
	//}

	srv.writeResponseBody(w, http.StatusCreated, []byte{})
}
