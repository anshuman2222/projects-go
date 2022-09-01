package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"io/ioutil"
	"net/http"
	"parking"
	"vehicle"

	"github.com/gorilla/mux"
)

var capa = &parking.Capacity{
	Total_four_wheeler_space: 3,
	Total_two_wheeler_space:  3,
}

var park = &parking.Parking{
	Capacity: capa,
}

func reserve_for_vehicle(w http.ResponseWriter, r *http.Request) {

	read, _ := ioutil.ReadAll(r.Body)

	var jsonMap vehicle.Vehicle
	json.Unmarshal([]byte(string(read)), &jsonMap)

	vehicle_no := jsonMap.Vehicle_number
	vehicle_type := jsonMap.Type

	vehicle_obj := park.Get_reserve_parking_space(vehicle_no, vehicle_type)

	park.Parking_info[vehicle_obj.Ticket.Parking_space] = vehicle_obj
	fmt.Println(park, park.Capacity, park.Parking_info, vehicle_obj)

}
func main() {

	r := mux.NewRouter()

	park.Parking_info = make(map[int]*vehicle.Vehicle)

	r.HandleFunc("/api/entry", reserve_for_vehicle).Methods("POST")

	http.ListenAndServe(":8000", r)
}
