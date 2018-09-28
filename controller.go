package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// SERVER the DB server
const SERVER = "mongodb://arcworker:arcworker123@ds111113.mlab.com:11113/arcworker"

// DBNAME the name of the DB instance
const DBNAME = "arcworker"

// COLLECTION is the name of the collection in DB
const COLLECTION = "worker"

// Index Is page default /
func Index(w http.ResponseWriter, r *http.Request) {
	results := "{'data': Running API ArcWork}"
	data, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// GetAll List All Workers
func GetAll(w http.ResponseWriter, r *http.Request) {
	//workers := Repository.GetWorkers() // list of all workers
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	fmt.Println("CONECTOU COM SUCESSO")

	//defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Workers{}

	fmt.Println("VARIAVEL C")

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	fmt.Println("RESULTADOS")

	// log.Println(workers)
	data, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddWorker POST /
func AddWorker(w http.ResponseWriter, r *http.Request) {
	var worker Worker
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request

	log.Println(body)

	if err != nil {
		log.Fatalln("Error AddWorker", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddWorker", err)
	}

	if err := json.Unmarshal(body, &worker); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddWorker unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//log.Println(worker)

	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// workerId += 1
	// worker.ID = workerId
	session.DB(DBNAME).C(COLLECTION).Insert(worker)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Added New Worker ID- ", worker.ID)

	//success := Repository.AddWorker(worker) // adds the worker to the DB

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}
