package users

import (
	shared "golang-mvc-boilerplate/server/sharedVariables"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Message is used to pass messages among model, view and controller
type Message struct {
	UserID  string        `json:"user_id,omitempty"`
	Status  int32         `json:"status"`
	Message string        `json:"message,omitempty"`
	User    *shared.User  `json:"user,omitempty"`
	Users   *shared.Users `json:"users,omitempty"`
	JWT     string        `json:"JWT,omitempty"`
}

//UpdateUser is used to represent a user while updating
type UpdateUser struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	Password    string `json:"new_password"`
}

//Result for isAdmin
type Result struct {
	Role string `bson:"role"`
}

//DataAccessLayer defines methods we need from the database
type DataAccessLayer interface {
	FindUser(collectionName string, userID string, userCheck *shared.User) error
	InsertUser(collectionName string, user shared.User) error
	UpdateUser(collectionName string, userID string, updatedUser shared.User) error
	GetUsers(collectionName string, userList *shared.Users) error
	DeleteUser(collectionName string, userID string) error
	IsAdmin(collectionName string, userID string, result *Result) error
	CloseSession()
	c(collectionName string) *mgo.Collection
}

// MongoSession is an implementation of DataAccessLayer for MongoDB
type MongoSession struct {
	session *mgo.Session
	dbName  string
}

var mgoSession MongoSession

//NewMongoSession is used to create or use previously created mongo sessions
func NewMongoSession(dbURI string, dbName string, dbUsername string, dbPassword string) (DataAccessLayer, Message) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{dbURI},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: dbUsername,
		Password: dbPassword,
	}
	if mgoSession.session == nil {
		var err error
		mgoSession.session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			log.Fatal(err)
			returnMessage := Message{
				Status:  500,
				Message: "Failed to establish connection to mongo server ",
			}
			mongo := &MongoSession{}
			return mongo, returnMessage
		}
	}
	returnMessage := Message{
		Status: 200,
	}
	mongo := &MongoSession{
		session: mgoSession.session.Copy(),
		dbName:  dbName,
	}
	return mongo, returnMessage
}

//CloseSession is used to close a mongodb session
func (m *MongoSession) CloseSession() {
	m.session.Close()
}

// c is a helper method to get a collection from the session
func (m *MongoSession) c(collectionName string) *mgo.Collection {
	return m.session.DB(m.dbName).C(collectionName)
}

// FindUser checks if a doc is present in db
func (m *MongoSession) FindUser(collectionName string, userID string, userCheck *shared.User) error {
	return m.c(collectionName).Find(bson.M{"userid": userID}).One(userCheck)
}

// InsertUser stores documents in mongo
func (m *MongoSession) InsertUser(collectionName string, user shared.User) error {
	return m.c(collectionName).Insert(user)
}

// UpdateUser updates documents in mongo
func (m *MongoSession) UpdateUser(collectionName string, userID string, updatedUser shared.User) error {
	return m.c(collectionName).Update(bson.M{"userid": userID}, updatedUser)
}

// GetUsers is used to get all the users from db
func (m *MongoSession) GetUsers(collectionName string, userList *shared.Users) error {
	return m.c(collectionName).Find(nil).Select(bson.M{"userid": 1, "role": 1}).All(userList)
}

// DeleteUser stores documents in mongo
func (m *MongoSession) DeleteUser(collectionName string, userID string) error {
	return m.c(collectionName).Remove(bson.M{"userid": userID})
}

// IsAdmin is used to check if a user is admin or not
func (m *MongoSession) IsAdmin(collectionName string, userID string, result *Result) error {
	return m.c(collectionName).Find(bson.M{"userid": userID}).Select(bson.M{"role": 1}).One(result)
}
