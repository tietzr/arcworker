package worker

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "http://localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "arcworker"

// COLLECTION is the name of the collection in DB
const COLLECTION = "worker"

// GetWorkers returns the list of Workers
func (r Repository) GetWorkers() Workers {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Workers{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddWorker adds a Worker in the DB
func (r Repository) AddWorker(worker Worker) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	workerId += 1
	worker.ID = workerId
	session.DB(DBNAME).C(COLLECTION).Insert(worker)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Worker ID- ", worker.ID)

	return true
}
