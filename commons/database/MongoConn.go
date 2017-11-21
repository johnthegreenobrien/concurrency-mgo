// This program provides a sample application for using MongoDB with
// the mgo driver.
package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"reflect"
	//"sync"
	"time"
)

/**
 * This file implements the singleton pattern to avoid the sistem to deal with more
 * than one Mongodb connection at once.
 */

const (
	MongoDBHosts = "localhost:27307"
	AuthDatabase = "delivery"
	AuthUserName = "deliveryUser"
	AuthPassword = "delivery123"
)

/**
 *
 */
type MongoSession struct {
	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

/**
 * @method db.sessione.Clone() GetMongoSession It create and instantiate a Mongodb connection
 * @return db.sessione.Clone()
 */
func ConnMongo() *MongoSession {
	db := &MongoSession{}
	return db.connect()
}

//var mgoSession *mgo.Session

// main is the entry point for the application.
func (db *MongoSession) connect() *MongoSession {

	if db.session == nil {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{MongoDBHosts},
			Timeout:  60 * time.Second,
			Database: AuthDatabase,
			Username: AuthUserName,
			Password: AuthPassword,
		}

		var err error
		db.session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			fmt.Println("CreateSession: %s\n", err)
		} else {
			fmt.Println("Session Created")
		}

		db.session.SetMode(mgo.Monotonic, true)
	}
	//db.session = db.session

	return db
}

func (db *MongoSession) GetSession() *mgo.Session {
	return db.session
}

func (db *MongoSession) GetIncrementer(field string) mgo.Change {
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{field: 1}},
		ReturnNew: true,
	}
	return change
}
