package reserve

import (
	lib "lib/common"
	"log"
	"sync"
)

type CustomerReserveDeviceResponse struct {
	CustomerId      string        `json:"customer_id"`
	ReservedDevices []*lib.Device `json:"reserved_devices"`
}

// Device channels initilisation
var Reserve_device_chan = make(chan lib.DevicesJson, 100)
var Device_res_chan = make(chan map[string]*lib.Device, 100)

// Customer device channel initilisation
var Customer_device_reserve_channel = make(chan lib.CutomersDeviceJson, 100)
var Customer_device_reserve_channel_res = make(chan []*CustomerReserveDeviceResponse, 100)

// This function is used to check whether devices are registered or not
func is_devices_exists(d_spec lib.DevicesJson) bool {
	var device_id string

	for i := 0; i < len(d_spec.DeviceList); i++ {
		device_id = d_spec.DeviceList[i].ID

		if _, ok := lib.Devices[device_id]; ok {
			log.Printf("Device %s is already configured", device_id)
			return true
		}
	}
	return false
}

// This function is used to register the devices
func Reserve_device(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer func() {
		lib.Stack_trace()
		wg.Done()
	}()

	for {
		select {
		case devices_json := <-Reserve_device_chan:
			device_res := make(map[string]*lib.Device, 0)

			is_device_exists := is_devices_exists(devices_json)

			if is_device_exists {
				Device_res_chan <- nil
			} else {
				for i := 0; i < len(devices_json.DeviceList); i++ {
					// Create device object
					device := new(lib.Device)
					device.ID = devices_json.DeviceList[i].ID
					device.GPSCoordinate = devices_json.DeviceList[i].GPSCoordinate
					device.IPAddress = devices_json.DeviceList[i].SerialNumber
					device.SerialNumber = devices_json.DeviceList[i].SerialNumber

					mutex.Lock()
					log.Printf("Adding device: %s with info %v", device.ID, device)
					lib.Devices[device.ID] = device
					mutex.Unlock()

					device_res[devices_json.DeviceList[i].ID] = lib.Devices[devices_json.DeviceList[i].ID]
				}
				Device_res_chan <- device_res
			}
		}
	}
}

// This function is used to check whether customer exists or not.
func is_customer_exists_from_customer_device_reserve_payload(customers_json lib.CutomersDeviceJson) bool {
	for i := 0; i < len(customers_json.CutomersDeviceList); i++ {
		customer_id := customers_json.CutomersDeviceList[i].ID

		if _, ok := lib.Customers[customer_id]; ok {
			log.Printf("Customer with ID %s is already exists", customer_id)
			return true
		}
	}
	return false
}

// This function is used to check whether devices are already registered to customers or not
func check_reserved_device_exists_from_customer_device_reserve_payload(customers_json lib.CutomersDeviceJson) bool {
	device_id_list := make(map[string]bool)

	for i := 0; i < len(customers_json.CutomersDeviceList); i++ {

		for dev_inx := 0; dev_inx < len(customers_json.CutomersDeviceList[i].Devices); dev_inx++ {
			device_id := customers_json.CutomersDeviceList[i].Devices[dev_inx]
			if _, ok := device_id_list[device_id]; ok {
				return true
			}
			// Check if device is registered or not, will move this logic to separate function in future
			if _, ok := lib.Devices[device_id]; !ok {
				return true
			}
			device_id_list[device_id] = true
		}
	}
	return false
}

// This function is used to reserve the devices to customers
func Customer_reserve_device_handler(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer lib.Stack_trace()

	for {
		select {
		case customers_json := <-Customer_device_reserve_channel:
			var customer_id, device_id string

			// Check whether customer is already exists and devices are not marked as reserved
			if is_customer_exists_from_customer_device_reserve_payload(customers_json) ||
				check_reserved_device_exists_from_customer_device_reserve_payload(customers_json) {
				Customer_device_reserve_channel_res <- nil
			} else {
				customer_device_reserve_res := make([]*CustomerReserveDeviceResponse, 0)

				for i := 0; i < len(customers_json.CutomersDeviceList); i++ {
					customer_id = customers_json.CutomersDeviceList[i].ID
					lib.Customers[customer_id] = make(map[string]*lib.Device)

					reserved_devices_res := make([]*lib.Device, 0)

					for dev_inx := 0; dev_inx < len(customers_json.CutomersDeviceList[i].Devices); dev_inx++ {

						device_id = customers_json.CutomersDeviceList[i].Devices[dev_inx]

						mutex.Lock()
						log.Printf("Adding device %s to customer %s", device_id, customer_id)
						lib.Customers[customer_id][device_id] = lib.Devices[device_id]

						// Change device state to claimed
						log.Printf("Marking device %s as reserved / claimed", device_id)
						lib.Devices_state.ResearvedDevices[device_id] = true

						if _, ok := lib.Devices_state.UnresearvedDevices[device_id]; ok {
							delete(lib.Devices_state.UnresearvedDevices, device_id)
						}
						mutex.Unlock()

						reserved_devices_res = append(reserved_devices_res, lib.Devices[device_id])
					}
					//customer_device_reserve_res_obj = new(CustomerReserveDeviceResponse)
					var customer_device_reserve_res_obj = new(CustomerReserveDeviceResponse)
					customer_device_reserve_res_obj.CustomerId = customer_id
					customer_device_reserve_res_obj.ReservedDevices = make([]*lib.Device, 0)
					customer_device_reserve_res_obj.ReservedDevices = reserved_devices_res

					customer_device_reserve_res = append(customer_device_reserve_res, customer_device_reserve_res_obj)
				}
				log.Printf("Customers with device: %v", lib.Customers)
				Customer_device_reserve_channel_res <- customer_device_reserve_res
			}
		}
	}
}
