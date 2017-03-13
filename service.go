package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// NewServer initializes the service with the given Database, and sets up appropriate routes.
func NewServer(db *Database) *Server {
	router := httprouter.New()
	server := &Server{
		router: router,
		db:     db,
	}

	server.setupRoutes()
	return server
}

// Server contains all that is needed to respond to incoming requests, like a database. Other services like a mail,
// redis, or S3 server could also be added.
type Server struct {
	router *httprouter.Router
	db     *Database
}

// The ServerError type allows errors to provide an appropriate HTTP status code and message. The Server checks for
// this interface when recovering from a panic inside a handler.
type ServerError interface {
	HttpStatusCode() int
	HttpStatusMessage() string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	s.router.POST("/contacts", s.AddContact)
	s.router.GET("/contacts/:email", s.GetContactByEmail)

	// By default the router will handle errors. But the service should always return JSON if possible, so these
	// custom handlers are added.

	s.router.NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			writeJSONError(w, http.StatusNotFound, "")
		},
	)

	s.router.HandleMethodNotAllowed = true
	s.router.MethodNotAllowed = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			writeJSONError(w, http.StatusMethodNotAllowed, "")
		},
	)

	s.router.PanicHandler = func(w http.ResponseWriter, r *http.Request, e interface{}) {
		serverError, ok := e.(ServerError)
		if ok {
			writeJSONError(w, serverError.HttpStatusCode(), serverError.HttpStatusMessage())
		} else {
			log.Printf("Panic during request: %v", e)
			writeJSONError(w, http.StatusInternalServerError, "")
		}
	}
}

// AddContact handles HTTP requests to add a Contact.
func (s *Server) AddContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var contact Contact

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&contact); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Error decoding JSON")
		return
	}

	contactId, err := s.db.AddContact(contact)
	if err != nil {
		panic(err)
		return
	}
	contact.Id = contactId

	writeJSON(
		w,
		http.StatusCreated,
		&ContactResponse{
			Contact: &contact,
		},
	)
}

// AddContact handles HTTP requests to GET a Contact by an email address.
func (s *Server) GetContactByEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email, err := url.QueryUnescape(ps.ByName("email"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid email.")
		return
	}

	email = strings.TrimSpace(email)
	if email == "" {
		writeJSONError(w, http.StatusBadRequest, "Expected a single email.")
		return
	}

	contact, err := s.db.GetContactByEmail(email)
	if err != nil {
		writeUnexpectedError(w, err)
	} else if contact == nil {
		writeJSONNotFound(w)
	} else {
		writeJSON(
			w,
			http.StatusOK,
			&ContactResponse{
				Contact: contact,
			},
		)
	}
}

// ===== JSON HELPERS ==================================================================================================

func writeJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	if message == "" {
		message = http.StatusText(statusCode)
	}

	writeJSON(
		w,
		statusCode,
		&ErrorResponse{
			StatusCode: statusCode,
			Message:    message,
		},
	)
}

func writeJSONNotFound(w http.ResponseWriter) {
	writeJSONError(w, http.StatusNotFound, "")
}

func writeUnexpectedError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusInternalServerError, err.Error())
}
