package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aebranton/rest-api/internal/user"
	"github.com/gorilla/mux"
)

// Handler stores a pointer to our router and user service
type Handler struct {
	Router  *mux.Router
	Service *user.Service
}

// Response - simple struct for displaying results in json on a page if the request
// has no object to return (ie Delete requests)
type Response struct {
	Message string
	Error   string
}

// NewHandler - creates a new Handler
func NewHandler(service *user.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// InitRoutes - sets up the routes on our handler
func (h *Handler) InitRoutes() {
	fmt.Println("Building Routes")
	h.Router = mux.NewRouter()

	// Add user routes
	h.Router.HandleFunc("/api/user/{id}", h.GetUser).Methods("GET")
	h.Router.HandleFunc("/api/user", h.GetAllUsers).Methods("GET")
	h.Router.HandleFunc("/api/user", h.GetAllUsers).Queries("username", "{username}").Methods("GET")
	h.Router.HandleFunc("/api/user", h.CreateUser).Methods("POST")
	h.Router.HandleFunc("/api/user/{id}", h.DeleteUser).Methods("DELETE")
	h.Router.HandleFunc("/api/user/{id}", h.UpdateUser).Methods("PUT")

	// Adding a simple status check to make sure its online
	h.Router.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		response := Response{Message: "Status is okay!"}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			panic(err)
		}
	})
}

// WriteResponseMessage - helper for writing a response message to a page.
// Can be given any status code, and any message.
// Internally it will set the response pages header to the status code given, and then select
// wether the supplied message should go into the Message field, or the Error field, based on the code.
// Panics if anything goes wrong encoding the message to json.
func (h *Handler) WriteResponseMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	var response Response
	if status == http.StatusOK {
		response.Message = msg
	} else if status == http.StatusBadRequest {
		response.Error = msg
	} else {
		response.Message = msg
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}

// GetUintFromVars - Helper function - given a map which is taken from mux.Vars(reader),
// convert it to a uint. the string return is the original value taken from the map with the given
// key, in case the cast to uint fails - this allows us to log what was received and why it was an error
func (h *Handler) GetUintFromVars(vars map[string]string, key string) (uint, string, error) {
	k := vars[key]
	i, err := strconv.ParseUint(k, 10, 64)

	if err != nil {
		return 0, "", err
	}

	return uint(i), k, nil
}

// GetUser - gets a user given the ID from a query (.../user/1)
// Writes either the selected user as json and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, val, err := h.GetUintFromVars(vars, "id")

	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID given: %s", val))
		return
	}

	user, err := h.Service.GetUser(id)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Error getting user with ID: %d", id))
		return
	}

	user.ToJSON(w)
}

// GetUserByUsername - gets a user given the username from a query (.../user?username=test)
// Writes either the selected user as json and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		h.WriteResponseMessage(w, http.StatusBadRequest, "Invalid, or no username given")
		return
	}

	user, err := h.Service.GetUserByUsername(username)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Error getting user with username: %s", username))
		return
	}

	user.ToJSON(w)
}

// GetAllUsers - gets all users from the database.
// Not currently paginated or using limits for the purposes of this demo.
// Production solutions would allow something such as (../user?limit=20&offset=40)
// Writes either the selected user as json and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()

	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, "Unable to retreive users")
		return
	}

	users.ToJSON(w)
}

// CreateUser - adds a user to the database with the given information in the body as JSON
// Writes the created user as json and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, "Failed to decode user from requests JSON")
		return
	}

	pwd, err := user.HashPassword(newUser.Password)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, "Unable to create user - password failed to hash.")
		return
	}

	newUser.Password = pwd
	user, err := h.Service.CreateUser(newUser)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Unable to create new user: %s", err))
		return
	}

	user.ToJSON(w)
}

// UpdateUser - updates a user in the database with the given id, and updates the supplied fields/data
// in the request body as JSON
// Writes the updated user as json and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, val, err := h.GetUintFromVars(vars, "id")
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID given: %s", val))
		return
	}

	var updatedUser user.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, "Failed to decode user from requests JSON")
		return
	}

	user, err := h.Service.UpdateUser(id, updatedUser)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Unable to update user with ID: %d", id))
		return
	}

	user.ToJSON(w)
}

// DeleteUser - Deletes a user from the database with the given ID
// Writes a success Response message and a 200 status code
// or a 400 status code, and an error message written within a Response object
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, val, err := h.GetUintFromVars(vars, "id")
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Invalid user ID given: %s", val))
		return
	}

	err = h.Service.DeleteUser(id)
	if err != nil {
		h.WriteResponseMessage(w, http.StatusBadRequest, fmt.Sprintf("Unbale to delete user with ID: %d", id))
	}

	h.WriteResponseMessage(w, http.StatusOK, fmt.Sprintf("Success deleting comment: %d", id))
}
