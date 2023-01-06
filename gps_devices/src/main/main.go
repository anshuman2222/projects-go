package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	Fetch "fetch"
	lib "lib/common"
	Release "release"
	Reserve "reserve"

	"github.com/gorilla/mux"
)

// Global variable for mutex and waitGroup
var mutex sync.Mutex
var wg sync.WaitGroup

// This handler is used to register the devices
func reserve_devices_handler(w http.ResponseWriter, r *http.Request) {
	defer lib.Stack_trace()

	// Read the data from payload
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed in reading reserve devices payload. Error: %s", err.Error())
		return
	}
	var devices_json lib.DevicesJson

	// Unmarshal the payload
	err = json.Unmarshal(payload, &devices_json)

	if err != nil {
		log.Printf("Failed in unmarshal the payload. Error: %s", err.Error())
	}

	// Send the devices info to channel to store it
	Reserve.Reserve_device_chan <- devices_json

	// Get the status of the devices after storing it
	device_res := <-Reserve.Device_res_chan

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")

	var res_in_bytes []byte

	if device_res != nil {
		res_in_bytes, err = json.Marshal(device_res)
		if err != nil {
			log.Printf("Failed in marshaling data: %s . Error is: %s", device_res, err.Error())
		}
	} else {
		res_in_bytes = []byte("Devices are already configured.")
	}
	_, err = w.Write(res_in_bytes)

	if err != nil {
		log.Printf("Failed in writting response. Error: %s", err.Error())
	}
}

// This handler is used to reserved the devices to customers
func customer_reserve_device_handler(w http.ResponseWriter, r *http.Request) {
	defer lib.Stack_trace()

	// Read the data from payload
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed in reading customer reserve devices payload. Error: %s", err.Error())
	}
	var customers_json lib.CutomersDeviceJson

	err = json.Unmarshal(payload, &customers_json)

	if err != nil {
		log.Printf("Failed in unmarshal payload. Error: %s", err.Error())
	}
	// Send customer device info to channel to reserve the devices of customer
	Reserve.Customer_device_reserve_channel <- customers_json

	// Get the response after reserving the devices of customers
	customers_res := <-Reserve.Customer_device_reserve_channel_res

	// Send response back to client
	w.Header().Set("Content-Type", "application/json")
	var res_in_bytes []byte

	if customers_res != nil {
		res_in_bytes, err = json.Marshal(customers_res)

		if err != nil {
			log.Printf("Failed in marshaling data %s . Error: %s", customers_res, err.Error())
		}
	} else {
		res_in_bytes = []byte("Either customer already exists or Devices not found")
	}
	_, err = w.Write(res_in_bytes)

	if err != nil {
		log.Printf("Failed in writting response. Error: %s", err.Error())
	}
}

// This handler is used to release the devices from customers
func release_customer_devices_handler(w http.ResponseWriter, r *http.Request) {
	defer lib.Stack_trace()

	customer_id := strings.Split(r.RequestURI, "/")[4]
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Failed in reading data. Error: %s", err.Error())
	}

	log.Printf("Releaseing devices from customer with ID: %s", customer_id)

	var customer_device_release lib.CustomerDeviceRelease

	err = json.Unmarshal(payload, &customer_device_release)

	if err != nil {
		log.Printf("Failed in unmarshaling data. Error: %s", err.Error())
	}
	customer_device_release.CustomerId = customer_id

	// Send devices of customer to channel to release the devices
	Release.Release_customer_device_channel <- customer_device_release

	// Send response back to client
	w.Header().Set("Content-Type", "application/json")

	if <-Release.Release_device_channel_res {
		w.Write([]byte(fmt.Sprintf("Successfully released devices from customer %s", customer_id)))
	} else {
		w.Write([]byte(fmt.Sprintf("Failed in releasing customer %s devices", customer_id)))
	}
}

// This handler is used to get the state of devices (claimed / unclaimed)
func get_devices_state_handler(w http.ResponseWriter, r *http.Request) {
	defer lib.Stack_trace()

	log.Println("Getting device state data")
	Fetch.Device_state_channel <- true

	device_res := <-Fetch.Device_state_channel_res

	log.Printf("Device state data are: %v", device_res)

	// Send devices state to client
	w.Header().Set("Content-Type", "application/json")
	byte_data, err := json.Marshal(device_res)

	if err != nil {
		log.Printf("Unable to convert clamined / unclaimed data to json. Erro: %s", err.Error())
	}
	w.Write(byte_data)
}

// This handler is used to get all the reserved devices to customer
func get_customer_devices_handler(w http.ResponseWriter, r *http.Request) {
	customer_id := strings.Split(r.RequestURI, "/")[4]

	if _, ok := lib.Customers[customer_id]; !ok {
		w.Write([]byte(fmt.Sprintf("Customer with ID %s doesn't exists", customer_id)))
	} else {
		res_in_bytes, err := json.Marshal(lib.Customers[customer_id])

		if err != nil {
			fmt.Println("Failed in marhsaling data. Error is: %s", err.Error())
		}
		// Sent all devices of customer to client
		w.Header().Set("Content-Type", "application/json")

		w.Write(res_in_bytes)
	}
}

func init() {
	defer lib.Stack_trace()

	lib.Customers = make(map[string]map[string]*lib.Device)
	lib.Devices = make(map[string]*lib.Device)
	lib.Devices_state.ResearvedDevices = make(map[string]bool)
	lib.Devices_state.UnresearvedDevices = make(map[string]bool)

	wg.Add(4)
	go Reserve.Reserve_device(&wg, &mutex)
	go Release.Release_customer_device(&wg, &mutex)
	go Reserve.Customer_reserve_device_handler(&wg, &mutex)
	go Fetch.Get_device_status(&wg)
}

func main() {
	defer lib.Stack_trace()

	log.Println("Program is starting....")

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/devices", reserve_devices_handler).Methods("POST")
	r.HandleFunc("/api/v1/customer/reserve/devices", customer_reserve_device_handler).Methods("POST")
	r.HandleFunc("/api/v1/devices/state", get_devices_state_handler).Methods("GET")
	r.HandleFunc("/api/v1/customers/{id}/devices", release_customer_devices_handler).Methods("DELETE")
	r.HandleFunc("/api/v1/customers/{id}/devices", get_customer_devices_handler).Methods("GET")

	http.ListenAndServe(":8000", r)
	wg.Wait()
}
