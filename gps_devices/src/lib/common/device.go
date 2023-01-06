package common

import (
	"fmt"
	"log"
	"runtime/debug"
)

// Structure of device data
type Device struct {
	ID            string
	GPSCoordinate string
	SerialNumber  string
	IPAddress     string
}

// Store the devices in map
var Devices map[string]*Device

// Json encoding / deconding for device data
type DevicesJson struct {
	DeviceList []struct {
		ID            string `json:"id"`
		GPSCoordinate string `json:"gps_coordinate"`
		SerialNumber  string `json:"serial_number"`
		IPAddress     string `json:"ip_address"`
	} `json:"devices"`
}

// Store the customer data in map
var Customers map[string]map[string]*Device

// Store customer information with devices
type Customer struct {
	ID      string
	Devices map[string]*Device
}

// Json encoding / deconding for customer devices
type CutomersDeviceJson struct {
	CutomersDeviceList []struct {
		ID      string   `json:"id"`
		Devices []string `json:"devices"`
	} `json:"customers"`
}

var Devices_state DeviceState

// store device state info
type DeviceState struct {
	ResearvedDevices   map[string]bool
	UnresearvedDevices map[string]bool
}

// Json encoding / deconding for deleting devicing of customer
type CustomerDeviceRelease struct {
	CustomerId string
	Devices    []string `json:"devices"`
}

func Stack_trace() {
	if r := recover(); r != nil {
		errStr := fmt.Sprintf("Error Recovered %s\n with Panic traceback : %s", r, debug.Stack())
		log.Println(errStr)
	}
}
