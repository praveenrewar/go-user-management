package users

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//Controller ...
type Controller struct {
	Repository Repository
}

// GetUsers GET /
func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := c.Repository.GetUsers() // list of all users
	data, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// Signup POST /
func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	type Message struct {
		UserID  string `json:"user_id,omitempty"`
		Status  int32  `json:"status"`
		Message string `json:"message,omitempty"`
	}
	body, err := ioutil.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Fatalln("Error in signing up", err)
		w.WriteHeader(http.StatusInternalServerError)
		message := &Message{
			Status:  500,
			Message: "Error in signing up"}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error in signing up", err)
	}
	if err := json.Unmarshal(body, &user); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error Signup unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			message := &Message{
				Status:  500,
				Message: "Error Signup unmarshalling data"}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}
	}
	success := c.Repository.AddUser(user) // adds the user to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		message := &Message{
			Status:  500,
			Message: "Internal Server Error"}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	message := &Message{
		Status:  200,
		UserID:  user.UserID,
		Message: "Signup successfull"}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}

// UpdateProfile PUT /
func (c *Controller) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	if err != nil {
		log.Fatalln("Error UpdateUser", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdateUser", err)
	}
	if err := json.Unmarshal(body, &user); err != nil { // unmarshall body contents as a type Candidate
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateUser unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.Repository.UpdateUser(user) // updates the user in the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// DeleteUser DELETE /
func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] // param id

	if err := c.Repository.DeleteUser(id); err != "" { // delete a users by id
		if strings.Contains(err, "404") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err, "500") {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}
