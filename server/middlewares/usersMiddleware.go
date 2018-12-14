package usersmiddlewares

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func signupValidator(w http.ResponseWriter, r *http.Request) {
	type FormData struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
	var user FormData
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
	log.Println("User Data", user)
}
