package release

import (
	lib "lib/common"
	"log"
	"sync"
)

var Release_customer_device_channel = make(chan lib.CustomerDeviceRelease, 100)
var Release_device_channel = make(chan string, 100)
var Release_device_channel_res = make(chan bool, 100)

// This function is used to unregister the devices
func Release_device(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer func() {
		lib.Stack_trace()
		wg.Done()
	}()

	for {
		select {
		case device_id := <-Release_device_channel:
			if _, ok := lib.Devices[device_id]; !ok {
				log.Println("Device with ID %s doesn't exists", device_id)
				Release_device_channel_res <- false
			} else {
				mutex.Lock()
				log.Printf("Delete the device: %s", device_id)
				delete(lib.Devices, device_id)
				Release_device_channel_res <- true
				mutex.Unlock()
			}
		}
	}
}

// This function is used to release the devices from the customers
func Release_customer_device(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer func() {
		lib.Stack_trace()
		wg.Done()
	}()

	for {
		select {
		case payload := <-Release_customer_device_channel:
			customer_id := payload.CustomerId
			customer := lib.Customers[customer_id]

			for _, device_id := range payload.Devices {
				if _, ok := customer[device_id]; ok {
					mutex.Lock()
					log.Printf("Releasing the device %s from customer: %s", device_id, customer_id)
					delete(customer, device_id)
					delete(lib.Devices_state.ResearvedDevices, device_id)
					lib.Devices_state.UnresearvedDevices[device_id] = true
					mutex.Unlock()
				} else {
					log.Printf("Device %s doesn't registers with customer %s ", device_id, customer_id)
				}
			}
			// Currently returing true as operation is successful, Ignoring unregistered device ID for now
			Release_device_channel_res <- true
		}
	}
}
