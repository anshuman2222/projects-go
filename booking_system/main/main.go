package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"io/ioutil"
	"net/http"
	"parking"

	"github.com/gorilla/mux"
)

var capa = &parking.Capacity{Four_wheeler: 3, Two_wheeler: 3}
var park = &parking.Parking{Capacity: capa, Four_wheeler: 2, Two_wheeler: 3}

func reserve_for_vehicle(w http.ResponseWriter, r *http.Request) {

	read, err := ioutil.ReadAll(r.Body)

	fmt.Println(string(read), err)

	var jsonMap map[string]int
	json.Unmarshal([]byte(string(read)), &jsonMap)

	vehicle_no := jsonMap["vehicle_no"]

	is_parked := park.Get_reserve_parking_space(vehicle_no)

	fmt.Println(is_parked)
}
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/entry", reserve_for_vehicle).Methods("POST")

	http.ListenAndServe(":8000", r)
}
