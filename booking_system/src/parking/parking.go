package parking

import (
	"log"
	"vehicle"
)

type Capacity struct {
	Four_wheeler int
	Two_wheeler  int
}
type Parking struct {
	Capacity     *Capacity       `json:"capacity"`
	Four_wheeler int             `json:"four_wheeler`
	Two_wheeler  int             `json:"two_wheeler`
	Vehicle      vehicle.Vehicle `json:"vehicle"`
}

func (p *Parking) Is_parking_available_for_two_wheeler() bool {
	if p.Two_wheeler == 0 {
		return false
	} else {
		return true
	}
}

func (p *Parking) Is_parking_available_for_four_wheeler() bool {
	if p.Four_wheeler == 0 {
		return false
	} else {
		return true
	}
}

func (p *Parking) Is_parking_available(vehicle_type int) bool {
	if vehicle_type == 2 {
		return p.Is_parking_available_for_two_wheeler()
	} else {
		return p.Is_parking_available_for_four_wheeler()
	}
}

func (p *Parking) Get_reserve_parking_space(vehicle_type int) bool {
	is_parking_available := p.Is_parking_available(vehicle_type)

	if is_parking_available == false {
		log.Fatal("There is no parking space left of type: ", vehicle_type)
		return false
	}
	if vehicle_type == 2 {
		p.Two_wheeler -= 1
	} else {
		p.Four_wheeler -= 1
	}
	return true
}
