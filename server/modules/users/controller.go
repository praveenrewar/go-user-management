package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../../middlewares/jwtAuthenticate"
	"../../middlewares/usersMiddleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//Controller ...
type Controller struct {
	Repository Repository
}

// Signup POST /
func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {

	var user User
	data := context.Get(r, "user")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	user = User(data.(usersmiddleware.UserFormData))
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	user.Password = string(hash)
	result := c.Repository.Signup(user) // adds the user to the DB
	if result.Status != 201 {
		w.WriteHeader(int(result.Status))
		message := Message{
			Status:  result.Status,
			Message: result.Message}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}
	w.WriteHeader(http.StatusCreated)
	message := &Message{
		Status:  201,
		UserID:  user.UserID,
		Message: "Signup successfull"}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}

// Login POST /
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var user User
	data := context.Get(r, "user")
	user = User(data.(usersmiddleware.UserFormData))

	result := c.Repository.Login(user)
	if result.Status != 200 {
		w.WriteHeader(int(result.Status))
		message := Message{
			Status:  result.Status,
			Message: result.Message}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}
	plainPassword := []byte(user.Password)
	hashPassword := []byte(result.User.Password)
	comparePassword := bcrypt.CompareHashAndPassword(hashPassword, plainPassword)
	if comparePassword != nil {
		w.WriteHeader(401)
		message := Message{
			Status:  401,
			Message: "Invalid password"}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": result.User.UserID,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status:  200,
		UserID:  user.UserID,
		Message: "Login Successfull",
		JWT:     tokenString}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}

// UpdatePassword POST /
func (c *Controller) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var user UpdateUser
	data := context.Get(r, "user")
	user = UpdateUser(data.(usersmiddleware.UpdateProfileFormData))

	result := c.Repository.UpdateUser(user) // updates the user in the DB
	w.WriteHeader(int(result.Status))
	message := Message{
		Status:  result.Status,
		UserID:  user.UserID,
		Message: result.Message}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}

// GetUsers POST /
func (c *Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result := c.Repository.GetUsers() // get all the users from db
	w.WriteHeader(int(result.Status))
	message := Message{
		Status:  result.Status,
		Message: result.Message,
		Users:   result.Users}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}

// DeleteUser DELETE /
func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"] // param user_id
	JWTData := context.Get(r, "user_jwt")
	JWTUser := JWTData.(jwtauthenticate.UserJWTData)
	isAdmin := c.Repository.IsAdmin(JWTUser.UserID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(int(isAdmin.Status))
	if isAdmin.Status != 200 {
		message := Message{
			Status:  isAdmin.Status,
			Message: isAdmin.Message}
		userJSON, _ := json.Marshal(message)
		w.Write(userJSON)
		return
	}
	result := c.Repository.DeleteUser(userID) // delete a users by user_id
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(int(result.Status))
	message := Message{
		Status:  result.Status,
		Message: result.Message,
		UserID:  result.UserID}
	userJSON, _ := json.Marshal(message)
	w.Write(userJSON)
	return
}
