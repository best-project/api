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

	user, err := srv.db.User.GetByName(name)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrapf(err, "while getting user with name: %s", name))
		return
	}

	writeResponseObject(w, http.StatusOK, user)
}

func (srv *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	user := srv.getUserData(w, r)

	token, err := NewJWT(NewCustomPayload(user))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.New("while creating token"))
		return
	}

	srv.writeResponseBody(w, http.StatusOK, []byte(token))
}

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	user := srv.getUserData(w, r)

	if srv.db.User.Exist(user) {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("[username, password] params required"))
		return
	}
	var err error
	user.Password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "while hashing password"))
		return
	}

	srv.db.User.SaveUser(user)
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
	fmt.Println(feed)

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
