package users

import (
	"log"

	"../../sharedVariables"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var address = shared.Address.(string)

//Repository ...
type Repository struct{}

// DBNAME the name of the DB instance
var DBNAME = shared.DbName.(string)

// DOCNAME the name of the document
const DOCNAME = "users"

var mgoSession *mgo.Session

//GetMongoSession is used to create or use previously created mongo sessions
func GetMongoSession() (*mgo.Session, Message) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(address)
		if err != nil {
			log.Fatal(err)
			returnMessage := Message{
				Status:  500,
				Message: "Failed to establish connection to mongo server ",
			}
			return nil, returnMessage
		}
	}
	returnMessage := Message{
		Status: 200,
	}
	return mgoSession.Copy(), returnMessage
}

// Login returns the list of Users
func (r Repository) Login(user User) Message {
	var userCheck User
	session, err := GetMongoSession()
	if err.Status != 200 {
		return err
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	findErr := collection.Find(bson.M{"userid": user.UserID}).One(&userCheck)
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
	session, dbMessage := GetMongoSession()
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	var userCheck User
	findQuery := collection.Find(bson.M{"userid": user.UserID}).One(&userCheck)
	if findQuery != nil {
		if findQuery.Error() == "not found" {
			err := session.DB(DBNAME).C(DOCNAME).Insert(user)
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
	session, dbMessage := GetMongoSession()
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	var userCheck User
	var updatedUser User
	findQuery := collection.Find(bson.M{"userid": user.UserID}).One(&userCheck)
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
		updateError := collection.Update(bson.M{"userid": user.UserID}, updatedUser)
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
	session, dbMessage := GetMongoSession()
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	var userCheck Users
	findQuery := collection.Find(nil).Select(bson.M{"userid": 1, "role": 1}).All(&userCheck)
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
	if len(userCheck) >= 0 {
		returnMessage := Message{
			Status:  200,
			Message: "User List",
			Users:   &userCheck,
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
	session, dbMessage := GetMongoSession()
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	if err := collection.Remove(bson.M{"userid": userID}); err != nil {
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
	session, dbMessage := GetMongoSession()
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.Close()
	collection := session.DB(DBNAME).C(DOCNAME)
	var result struct {
		Role string `bson:"role"`
	}
	if err := collection.Find(bson.M{"userid": userID}).Select(bson.M{"role": 1}).One(&result); err != nil {
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
