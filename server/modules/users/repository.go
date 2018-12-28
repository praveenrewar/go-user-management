package users

import (
	"log"

	"../../sharedVariables"
	"golang.org/x/crypto/bcrypt"
)

var address = shared.Address.(string)

//Repository ...
type Repository struct{}

// DBNAME the name of the DB instance
var DBNAME = shared.DbName.(string)

// DOCNAME the name of the document
const DOCNAME = "users"

// Login returns the list of Users
func (r Repository) Login(user User) Message {
	var userCheck User
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	findErr := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findErr != nil {
		if findErr.Error() == "not found" {
			returnMessage := Message{
				Status:  404,
				Message: "No such user exists",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	} else {
		if userCheck == (User{}) {
			returnMessage := Message{
				Status:  404,
				Message: "No such user exists",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  200,
			Message: "User found",
			User:    &userCheck,
		}
		return returnMessage
	}
}

// Signup inserts a user in the DB
func (r Repository) Signup(user User) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userCheck User
	findQuery := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			err := session.InsertUser(DOCNAME, user)
			if err != nil {
				log.Fatal(err)
				returnMessage := Message{
					Status:  500,
					Message: "Internal Server Error",
				}
				return returnMessage
			}
			returnMessage := Message{
				Status:  201,
				Message: "Signup successfull",
				UserID:  user.UserID,
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	}
	if userCheck != (User{}) {
		returnMessage := Message{
			Status:  409,
			Message: "UserID already exists",
		}
		return returnMessage
	}
	returnMessage := Message{
		Status:  500,
		Message: "Internal Server Error",
	}
	return returnMessage

}

// UpdateUser updates an User in the DB
func (r Repository) UpdateUser(user UpdateUser) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userCheck User
	var updatedUser User
	findQuery := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			returnMessage := Message{
				Status:  404,
				Message: "No such user exists",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	}
	if userCheck != (User{}) {
		plainPassword := []byte(user.OldPassword)
		hashPassword := []byte(userCheck.Password)
		comparePassword := bcrypt.CompareHashAndPassword(hashPassword, plainPassword)
		if comparePassword != nil {
			returnMessage := Message{
				Status:  401,
				Message: "Old password is invalid",
			}
			return returnMessage
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}
		user.Password = string(hash)
		updatedUser.UserID = user.UserID
		updatedUser.Password = user.Password
		updateError := session.UpdateUser(DOCNAME, user.UserID, updatedUser)
		if updateError != nil {
			returnMessage := Message{
				Status:  500,
				Message: "Some error while changing password\n" + string(updateError.Error()),
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  200,
			Message: "Password changed successfully",
		}
		return returnMessage
	}
	returnMessage := Message{
		Status:  500,
		Message: "Internal Server Error",
	}
	return returnMessage
}

//GetUsers is used to get all the users from the db
func (r Repository) GetUsers() Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userList Users
	findQuery := session.GetUsers(DOCNAME, &userList)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			returnMessage := Message{
				Status:  200,
				Message: "No user found",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	}
	if len(userList) >= 0 {
		returnMessage := Message{
			Status:  200,
			Message: "User List",
			Users:   &userList,
		}
		return returnMessage
	}
	returnMessage := Message{
		Status:  500,
		Message: "Internal Server Error",
	}
	return returnMessage
}

// DeleteUser deletes an User (not used for now)
func (r Repository) DeleteUser(userID string) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	if err := session.DeleteUser(DOCNAME, userID); err != nil {
		if err.Error() == "not found" {
			returnMessage := Message{
				Status:  404,
				Message: "No such user exists",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	}

	returnMessage := Message{
		UserID:  userID,
		Status:  200,
		Message: "User Deleted",
	}
	return returnMessage
}

//IsAdmin checks if user is admin or not
func (r Repository) IsAdmin(userID string) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var result Result
	if err := session.IsAdmin(DOCNAME, userID, &result); err != nil {
		if err.Error() == "not found" {
			returnMessage := Message{
				Status:  401,
				Message: "Invalid JWT credentials",
			}
			return returnMessage
		}
		returnMessage := Message{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	}
	if result.Role == "admin" {
		returnMessage := Message{
			UserID:  userID,
			Status:  200,
			Message: "User is an admin",
		}
		return returnMessage
	}
	returnMessage := Message{
		UserID:  userID,
		Status:  401,
		Message: "User is not an admin",
	}
	return returnMessage

}
