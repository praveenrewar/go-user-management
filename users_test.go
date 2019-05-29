package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "go-user-management/server/modules/users"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
)

var a App
var address interface{}
var dbName interface{}

type User struct {
	UserID   string
	Password string
	Role     string
}

func TestMain(m *testing.M) {
	a = App{}
	a.Run()
	address = a.DBAddress
	dbName = a.DBName
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearCollection() {
	session, err := mgo.Dial(address.(string))
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	session.DB(dbName.(string)).C("users").RemoveAll(nil)
}

func addUser(userID string, password string, role string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	password = string(hash)
	session, err := mgo.Dial(address.(string))
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	var user User
	user.UserID = userID
	user.Password = password
	user.Role = role
	session.DB(dbName.(string)).C("users").Insert(user)
}

func getToken(userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": userID,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	return tokenString
}

func TestSignup(t *testing.T) {
	clearCollection()
	payload := []byte(`{"user_id":"test_user","role":"admin","password":"password"}`)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["user_id"] != "test_user" {
		t.Errorf("Expected user_id to be 'test_user'. Got '%v'", m["name"])
	}
	if m["message"] != "Signup successfull" {
		t.Errorf("Expected message to be 'Signup successfull'. Got '%v'", m["message"])
	}
}

func TestSignupBadRequest(t *testing.T) {
	clearCollection()
	payload := []byte(`{"role":"admin","password":"password"}`)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
}

func TestSignupExistingUser(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	payload := []byte(`{"user_id":"test_user","role": "admin","password":"password"}`)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusConflict, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["message"] != "UserID already exists" {
		t.Errorf("Expected message to be 'UserID already exists'. Got '%v'", m["message"])
	}
}
func TestLoginNonExistingUser(t *testing.T) {
	clearCollection()
	payload := []byte(`{"user_id":"test_user","password":"password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "No such user exists" {
		t.Errorf("Expected the 'message' key of the response to be set to 'No such user exists'. Got '%s'", m["message"])
	}
}

func TestLogin(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	payload := []byte(`{"user_id":"test_user","password":"password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["user_id"] != "test_user" {
		t.Errorf("Expected user_id to be 'test_user'. Got '%v'", m["user_id"])
	}
	if m["message"] != "Login Successfull" {
		t.Errorf("Expected message to be 'Login Successfull'. Got '%v'", m["message"])
	}
}

func TestLoginBadRequest(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	payload := []byte(`{"user_id":"test_user"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLoginWrongPassword(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	payload := []byte(`{"user_id":"test_user","password":"wrong_password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Invalid password" {
		t.Errorf("Expected message to be 'Invalid password'. Got '%v'", m["message"])
	}
}

func TestUpdatePassword(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	payload := []byte(`{"old_password":"password","new_password":"password123"}`)
	req, _ := http.NewRequest("POST", "/update-password", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Password changed successfully" {
		t.Errorf("Expected message to be 'Password changed successfully'. Got '%v'", m["message"])
	}
}

func TestUpdatePasswordBadRequest(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	payload := []byte(`{"old_password":"password"}`)
	req, _ := http.NewRequest("POST", "/update-password", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUpdatePasswordSamePassword(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	payload := []byte(`{"old_password":"password","new_password":"password"}`)
	req, _ := http.NewRequest("POST", "/update-password", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "New password cannot be same as old password" {
		t.Errorf("Expected message to be 'New password cannot be same as old password'. Got '%v'", m["message"])
	}
}

func TestUpdatePasswordInvalidToken(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	// tokreturn nil, fmt.Errorf("There was an error")en := "token"
	payload := []byte(`{"old_password":"password","new_password":"password123"}`)
	req, _ := http.NewRequest("POST", "/update-password", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+"token")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Invalid authorization token" {
		t.Errorf("Expected message to be 'Invalid authorization token'. Got '%v'", m["message"])
	}
}

func TestUpdatePasswordInvalidOldPassword(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	payload := []byte(`{"old_password":"invalid_password","new_password":"password123"}`)
	req, _ := http.NewRequest("POST", "/update-password", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Old password is invalid" {
		t.Errorf("Expected message to be 'Old password is invalid'. Got '%v'", m["message"])
	}
}

func TestGetUsers(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	req, _ := http.NewRequest("GET", "/get-users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "User List" {
		t.Errorf("Expected message to be 'User List'. Got '%v'", m["message"])
	}
}

func TestGetUsersInvalidToken(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	req, _ := http.NewRequest("GET", "/get-users", nil)
	req.Header.Set("Authorization", "Bearer "+token+"e")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Invalid authorization token" {
		t.Errorf("Expected message to be 'Invalid authorization token'. Got '%v'", m["message"])
	}
}

func TestGetUsersInvalidTokenUserDeleted(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	token := getToken("test_user")
	clearCollection()
	req, _ := http.NewRequest("GET", "/get-users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Invalid authorization token" {
		t.Errorf("Expected message to be 'Invalid authorization token'. Got '%v'", m["message"])
	}
}

func TestGetUsersInvalidTokenFormat(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	req, _ := http.NewRequest("GET", "/get-users", nil)
	req.Header.Set("Authorization", "Bearer ")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "Please insert authorization header in the format 'Bearer {token}'" {
		t.Errorf("Expected message to be 'Please insert authorization header in the format 'Bearer {token}''. Got '%v'", m["message"])
	}
}

func TestGetUsersMissingAuthorizationHeader(t *testing.T) {
	clearCollection()
	addUser("test_user", "password", "admin")
	req, _ := http.NewRequest("GET", "/get-users", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "An authorization header is required" {
		t.Errorf("Expected message to be 'An authorization header is required''. Got '%v'", m["message"])
	}
}

func TestDeleteUser(t *testing.T) {
	clearCollection()
	addUser("admin_user", "password", "admin")
	addUser("test_user", "password", "datascientist")
	token := getToken("admin_user")
	req, _ := http.NewRequest("DELETE", "/delete-user/test_user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "User Deleted" {
		t.Errorf("Expected message to be 'User Deleted'. Got '%v'", m["message"])
	}
	if m["user_id"] != "test_user" {
		t.Errorf("Expected user_id to be 'test_user'. Got '%v'", m["user_id"])
	}
}

func TestDeleteUserNotAdmin(t *testing.T) {
	clearCollection()
	addUser("admin_user", "password", "datascientist")
	addUser("test_user", "password", "datascientist")
	token := getToken("admin_user")
	req, _ := http.NewRequest("DELETE", "/delete-user/test_user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["message"] != "User is not an admin" {
		t.Errorf("Expected message to be 'User Deleted'. Got '%v'", m["message"])
	}
}
