package server

import (
	"encoding/json"
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (srv *Server) getUserData(w http.ResponseWriter, r *http.Request) *internal.User {
	userDTO := &internal.UserDTO{}
	err := json.NewDecoder(r.Body).Decode(&userDTO)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while decoding body"))
		return nil
	}
	if len(userDTO.Username) == 0 || len(userDTO.Password) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("[username, password] params required"))
		return nil
	}
	return srv.converter.ToModel(userDTO)
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
	user.Password = []byte{}

	writeResponseObject(w, http.StatusOK, user)
}

func (srv *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	userData := srv.getUserData(w, r)

	users, err := srv.db.User.GetByName(userData.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting by name %s", userData.Username))
		return
	}
	if len(users) == 0 {
		writeErrorResponse(w, http.StatusNotFound, errors.New("user not found"))
		return
	}
	user := users[0]
	if err := bcrypt.CompareHashAndPassword(user.Password, userData.Password); err != nil {
		writeErrorResponse(w, http.StatusForbidden, errors.Wrap(err, "while comparing passwords"))
		return
	}

	token, err := NewJWT(NewCustomPayload(&user))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while creating token"))
		return
	}
	user.Token = token
	if err := srv.db.User.UpdateUser(&user); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while saving user"))
		return
	}
	user.Password = []byte{}

	writeResponseObject(w, http.StatusOK, user)
}

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	user := srv.getUserData(w, r)

	if srv.db.User.Exist(user) {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("user already exists"))
		return
	}
	var err error
	user.Password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while hashing password"))
		return
	}
	if err := srv.db.User.SaveUser(internal.NewUser(user)); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while saving user"))
		return
	}
	srv.writeResponseCode(w, http.StatusCreated)
}

func (srv *Server) redirectToFb(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://www.facebook.com/v2.10/dialog/oauth?client_id=%s&redirect_uri=%s", srv.fb.GetAppID(), srv.fbCallbackURL), http.StatusTemporaryRedirect)
}

func (srv *Server) createUserFb(w http.ResponseWriter, r *http.Request) {
	user := &internal.User{}

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

	token, err := NewJWT(NewCustomPayload(user))
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
