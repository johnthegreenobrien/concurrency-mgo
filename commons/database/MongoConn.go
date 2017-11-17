// This program provides a sample application for using MongoDB with
// the mgo driver.
package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"reflect"
	"sync"
	"time"
	"userv/modules/delivery/models"
)

/**
 * This file implements the singleton pattern to avoid the sistem to deal with more
 * than one Mongodb connection at once.
 */

/**
 *
 */
type mongoSession struct {
	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

/**
 * @method db.sessione.Clone() GetmongoSession It create and instantiate a Mongodb connection
 * @return db.sessione.Clone()
 */
func ConnMongo() *mongoSession {
	db := &mongoSession{}
	return db.connectSingleton()
}

const (
	MongoDBHosts = "localhost:27307"
	AuthDatabase = "delivery"
	AuthUserName = "deliveryUser"
	AuthPassword = "delivery123"
	TestDatabase = "delivery"
)

// main is the entry point for the application.
func (db *mongoSession) connectSingleton() *mongoSession {

	if db.session == nil {
		// We need this object to establish a session to our MongoDB.
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{MongoDBHosts},
			Timeout:  60 * time.Second,
			Database: AuthDatabase,
			Username: AuthUserName,
			Password: AuthPassword,
		}

		// Create a session which maintains a pool of socket connections
		// to our MongoDB.
		var err error
		db.session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			fmt.Println("CreateSession: %s\n", err)
		} else {
			fmt.Println("Session Created")
		}

		// Reads may not be entirely up-to-date, but they will always see the
		// history of changes moving forward, the data read will be consistent
		// across sequential queries in the same session, and modifications made
		// within the session will be observed in following queries (read-your-writes).
		// http://godoc.org/labix.org/v2/mgo#Session.SetMode
		db.session.SetMode(mgo.Monotonic, true)
	}
	//db.session = db.session

	return db
}

// RunQuery is a function that is launched as a goroutine to perform
// the MongoDB work.
//func (db *mongoSession) RunQuery(waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
func (db *mongoSession) RunQuery(waitGroup *sync.WaitGroup) {
	// Decrement the wait group count so the program knows this
	// has been completed once the goroutine exits.
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := db.session.Copy()
	//sessionCopy := db.Session.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(TestDatabase).C("ipPorts")

	//fmt.Printf("RunQuery : %d : Executing\n")

	// Retrieve the list of stations.
	var deliveries []models.Delivery
	err := collection.Find(nil).All(&deliveries)
	if err != nil {
		fmt.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	//fmt.Printf("RunQuery : %d : Count[%d]\n", len(deliveries))
	fmt.Println("Count: ", len(deliveries))
}