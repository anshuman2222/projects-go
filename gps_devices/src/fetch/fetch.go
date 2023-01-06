package fetch

import (
	lib "lib/common"
	"sync"
)

// Channels initilisation
var Device_state_channel = make(chan bool, 100)
var Device_state_channel_res = make(chan map[string][]*lib.Device, 100)

// This function is used to return claimed / unclaimed devices.
func Get_device_status(wg *sync.WaitGroup) {
	defer func() {
		lib.Stack_trace()
		wg.Done()
	}()

	for {
		select {
		case <-Device_state_channel:

			var claimed_devices_list = make([]*lib.Device, 0)
			var unclaimed_devices_list = make([]*lib.Device, 0)
			var device_state = make(map[string][]*lib.Device, 0)

			// Get claimed devices
			for key, _ := range lib.Devices_state.ResearvedDevices {
				claimed_devices_list = append(claimed_devices_list, lib.Devices[key])
			}
			device_state["claimed"] = claimed_devices_list

			// Get unclaimed devices
			for key, _ := range lib.Devices_state.UnresearvedDevices {
				unclaimed_devices_list = append(unclaimed_devices_list, lib.Devices[key])
			}
			device_state["unclaimed"] = unclaimed_devices_list

			Device_state_channel_res <- device_state
		}
	}
}
