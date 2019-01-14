package users

import (
	"log"

	shared "golang-mvc-boilerplate/server/sharedVariables"

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
func (r Repository) Login(user shared.User) Message {
	var userCheck shared.User
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	findErr := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findErr != nil {
		if findErr.Error() == "not found" {
			return noSuchUserMessage()
		}
		return internalServerErrorMessage()
	} else {
		if userCheck == (shared.User{}) {
			return noSuchUserMessage()
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
func (r Repository) Signup(user shared.User) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userCheck shared.User
	findQuery := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			err := session.InsertUser(DOCNAME, user)
			if err != nil {
				log.Fatal(err)
				return internalServerErrorMessage()
			}
			returnMessage := Message{
				Status:  201,
				Message: "Signup successfull",
				UserID:  user.UserID,
			}
			return returnMessage
		}
		return internalServerErrorMessage()
	}
	if userCheck != (shared.User{}) {
		returnMessage := Message{
			Status:  409,
			Message: "UserID already exists",
		}
		return returnMessage
	}
	return internalServerErrorMessage()

}

// UpdatePassword updates an User in the DB
func (r Repository) UpdatePassword(user UpdateUser) Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userCheck shared.User
	var updatedUser shared.User
	findQuery := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			return noSuchUserMessage()
		}
		return internalServerErrorMessage()
	}
	if userCheck != (shared.User{}) {
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
			return internalServerErrorMessage()
		}
		returnMessage := Message{
			Status:  200,
			Message: "Password changed successfully",
		}
		return returnMessage
	}
	return internalServerErrorMessage()
}

//GetUsers is used to get all the users from the db
func (r Repository) GetUsers() Message {
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	var userList shared.Users
	findQuery := session.GetUsers(DOCNAME, &userList)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			returnMessage := Message{
				Status:  200,
				Message: "No user found",
			}
			return returnMessage
		}
		return internalServerErrorMessage()
	}
	if len(userList) >= 0 {
		returnMessage := Message{
			Status:  200,
			Message: "User List",
			Users:   &userList,
		}
		return returnMessage
	}
	return internalServerErrorMessage()
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
			return noSuchUserMessage()
		}
		return internalServerErrorMessage()
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
		return internalServerErrorMessage()
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

func noSuchUserMessage() Message {
	returnMessage := Message{
		Status:  404,
		Message: "No such user exists",
	}
	return returnMessage
}

func internalServerErrorMessage() Message {
	returnMessage := Message{
		Status:  500,
		Message: "Internal Server Error",
	}
	return returnMessage
}
