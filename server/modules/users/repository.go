package users

import (
	"fmt"
	"log"

	"../../sharedVariables"
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

// GetUsers returns the list of Users
func (r Repository) GetUsers() Users {
	session, err := mgo.Dial(address)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(DOCNAME)
	results := Users{}
	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddUser inserts a user in the DB
func (r Repository) AddUser(user User) bool {
	session, err := mgo.Dial(address)
	defer session.Close()
	session.DB(DBNAME).C(DOCNAME).Insert(user)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// UpdateUser updates an User in the DB (not used for now)
func (r Repository) UpdateUser(user User) bool {
	session, err := mgo.Dial(address)
	defer session.Close()
	session.DB(DBNAME).C(DOCNAME).UpdateId(user.UserID, user)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// DeleteUser deletes an User (not used for now)
func (r Repository) DeleteUser(id string) string {
	session, err := mgo.Dial(address)
	defer session.Close()

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return "NOT FOUND"
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err = session.DB(DBNAME).C(DOCNAME).RemoveId(oid); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	// Write status
	return "OK"
}
