package jwtauthenticate

import (
	shared "go-user-management/server/sharedVariables"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//UserJWTData is used to get user_id from jwt
type UserJWTData struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	IAT    string `json:"IAT"`
}

//AuthenticationMessage is used to return message from jwtAuthenticate middleware
type AuthenticationMessage struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}

//DataAccessLayer defines methods we need from the database
type DataAccessLayer interface {
	FindUser(collectionName string, userID string, userCheck *shared.User) error
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
func NewMongoSession(dbURI string, dbName string) (DataAccessLayer, AuthenticationMessage) {
	if mgoSession.session == nil {
		var err error
		mgoSession.session, err = mgo.Dial(dbURI)
		if err != nil {
			log.Fatal(err)
			returnMessage := AuthenticationMessage{
				Status:  500,
				Message: "Failed to establish connection to mongo server ",
			}
			mongo := &MongoSession{}
			return mongo, returnMessage
		}
	}
	returnMessage := AuthenticationMessage{
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
