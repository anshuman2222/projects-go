package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Candidate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func create_entity(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var candidate Candidate

	err = json.Unmarshal(data, &candidate)

	if err != nil {
		log.Fatal(err)
	}
	candidates[candidate.Id] = &candidate
	fmt.Println(candidates)
}

func get_entity(w http.ResponseWriter, r *http.Request) {
	data := r.RequestURI

	id := strings.Split(data, "/")[2]

	fmt.Println(data, string(data), id, candidates[id])
}

var candidates map[string]*Candidate

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	candidates = make(map[string]*Candidate)

	r := mux.NewRouter()

	r.HandleFunc("/create/entity", create_entity).Methods("POST")
	r.HandleFunc("/get/{id}", get_entity).Methods("GET")
	//r.HandleFunc("/delete/{}", delete_entity).Methods("GET")
	// r.HandleFunc("/update/{}", get_entity).Methods("GET")

	http.ListenAndServe(":8000", r)
}
