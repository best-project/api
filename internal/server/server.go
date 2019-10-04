package server

import (
	"github.com/urfave/negroni"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"strconv"
)

type Server struct {
	logger logrus.FieldLogger
}

func NewServer() *Server {
	return &Server{
		logger: logrus.New(),
	}
}

// Handle creates an http handler
func (srv *Server) Handle() http.Handler {
	var rtr = mux.NewRouter()

	rtr.Path("/user/create").Methods(http.MethodPost).Handler(negroni.New(negroni.WrapFunc(srv.createUser)))
	rtr.Path("/user").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(srv.getUser)))

	rtr.Path("/status").Methods(http.MethodGet).Handler(negroni.New(negroni.WrapFunc(
		func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200"))
	})))

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(rtr)
	return n
}

func (srv *Server) writeResponseCode(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(strconv.Itoa(code)))
}

func (srv *Server) writeJSONResponse(w http.ResponseWriter, code int, object interface{}) {
	writeResponse(w, code, object)
}

func writeResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func (srv *Server) writeErrorResponse(w http.ResponseWriter, code int, errorMsg, desc string) {
	if srv.logger != nil {
		srv.logger.Warnf("Server responds with error: [HTTP %d]: [%s] [%s]", code, errorMsg, desc)
	}
	writeErrorResponse(w, code, errorMsg, desc)
}

// writeErrorResponse writes error response compatible with OpenServiceBroker API specification.
func writeErrorResponse(w http.ResponseWriter, code int, errorMsg, desc string) {
	dto := struct {
		// Error is a machine readable info on an error.
		// As of 2.13 Open Broker API spec it's NOT passed to entity querying the catalog.
		Error string `json:"error,optional"`

		// Desc is a meaningful error message explaining why the request failed.
		// see: https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md#broker-errors
		Desc string `json:"description,optional"`
	}{}

	if errorMsg != "" {
		dto.Error = errorMsg
	}

	if desc != "" {
		dto.Desc = desc
	}
	writeResponse(w, code, &dto)
}

//func httpBodyToDTO(r *http.Request, object interface{}) error {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		return err
//	}
//	defer r.Body.Close()
//
//	err = json.Unmarshal(body, object)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
