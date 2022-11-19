package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Candaidate struct {
	Id   string         `json:"id"`
	Name string         `json:"name"`
	Fee  map[string]int `json:"fee"`
}

func create_entity(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}
	var candaidate Candaidate

	err = json.Unmarshal(data, &candaidate)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")

	if _, ok := candidates[candaidate.Id]; !ok {
		candidates[candaidate.Id] = &candaidate
	} else {
		err_msg := "Key is already exists: " + candaidate.Id
		w.Write([]byte(err_msg))
		return
	}

	status, err := json.Marshal(candaidate)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(status)
}

func get_entity(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]

	if _, ok := candidates[id]; ok {
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(candidates[id])

		if err != nil {
			log.Fatal(err)
		} else {
			w.Write(data)
		}
	} else {
		err_msg := "Key doesn't exists: " + id
		w.Write([]byte(err_msg))
	}
}

var candidates map[string]*Candaidate

func main() {
	candidates = make(map[string]*Candaidate)

	r := mux.NewRouter()

	r.HandleFunc("/candidates", create_entity).Methods("POST")
	r.HandleFunc("/candidates/{id}", get_entity).Methods("GET")
	/*
		r.HandleFunc("/candidates/{id}", delete_entity).Methods("DELETE")
		r.HandleFunc("/candidates/{id}", update_entity).Methods("PUT")
		r.HandleFunc("/candidates/list", get_all_entities).Methods("POST")
	*/

	http.ListenAndServe(":8000", r)
}
