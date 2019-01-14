package jwtauthenticate

import (
	shared "golang-mvc-boilerplate/server/sharedVariables"
)

var address = shared.Address.(string)

//Repository ...
type Repository struct{}

// DBNAME the name of the DB instance
var DBNAME = shared.DbName.(string)

// DOCNAME the name of the document
const DOCNAME = "users"

//ValidateUser is used to validate a user from db
func ValidateUser(user UserJWTData) AuthenticationMessage {
	var userCheck shared.User
	session, dbMessage := NewMongoSession(address, DBNAME)
	if dbMessage.Status != 200 {
		return dbMessage
	}
	defer session.CloseSession()
	findErr := session.FindUser(DOCNAME, user.UserID, &userCheck)
	if findErr != nil {
		if findErr.Error() == "not found" {
			returnMessage := AuthenticationMessage{
				Status:  401,
				Message: "Invalid authorization token",
			}
			return returnMessage
		}
		returnMessage := AuthenticationMessage{
			Status:  500,
			Message: "Internal Server Error",
		}
		return returnMessage
	} else {
		if userCheck == (shared.User{}) {
			returnMessage := AuthenticationMessage{
				Status:  401,
				Message: "Invalid authorization token",
			}
			return returnMessage
		}
		returnMessage := AuthenticationMessage{
			Status: 200,
		}
		return returnMessage
	}
}
