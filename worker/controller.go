package worker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	workers := c.Repository.GetWorkers() // list of all workers
	// log.Println(workers)
	data, _ := json.Marshal(workers)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddWorker POST /
func (c *Controller) AddWorker(w http.ResponseWriter, r *http.Request) {
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

	log.Println(worker)
	success := c.Repository.AddWorker(worker) // adds the worker to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}
