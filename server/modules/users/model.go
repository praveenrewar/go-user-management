package users

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//User represents a user
type User struct {
	UserID   string `json:"user_id"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

//Message is used to pass messages among model, view and controller
type Message struct {
	UserID  string `json:"user_id,omitempty"`
	Status  int32  `json:"status"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Users   *Users `json:"users,omitempty"`
	JWT     string `json:"JWT,omitempty"`
}

//UpdateUser is used to represent a user while updating
type UpdateUser struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	Password    string `json:"new_password"`
}

//Users is an array of User
type Users []User

var mgoSession *mgo.Session

//DataAccessLayer defines methods we need from the database
type DataAccessLayer interface {
}

// MongoSession is an implementation of DataAccessLayer for MongoDB
type MongoSession struct {
	session *mgo.Session
	dbName  string
}

//NewMongoSession is used to create or use previously created mongo sessions
func NewMongoSession(dbURI string, dbName string) (DataAccessLayer, Message) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(dbURI)
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
	mongo := &MongoSession{
		session: mgoSession.Copy(),
		dbName:  dbName,
	}
	return mongo, returnMessage
}

// c is a helper method to get a collection from the session
func (m *MongoSession) c(collection string) *mgo.Collection {
	return m.session.DB(m.dbName).C(collection)
}

// FindOne checks if a doc is present in db
func (m *MongoSession) FindOne(collectionName string, user User, userCheck User) error {
	return m.c(collectionName).Find(bson.M{"userid": user.UserID}).One(&userCheck)
}

// Insert stores documents in mongo
func (m *MongoSession) Insert(collectionName string, user User) error {
	return m.c(collectionName).Insert(user)
}
