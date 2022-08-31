package entry

import (
	"fmt"
	"log"
	parking "parking"
	"time"
	vehicle "vehicle"
)

type Entry struct {
	Vehicle    vehicle.Vehicle
	Entry_time time.Time
	Parking    parking.Parking
	No         int
}

func (e *Entry) enter_vehicle(vehicle_no string, vehicle_type int) {
	is_reserve := e.Parking.get_reserve_parking_space(vehicle_type)

	if is_reserve == true {
		fmt.Print("Parking is successful")
	} else {
		log.Fatal("There is no parking space left of type: ", vehicle_type)
	}
}
