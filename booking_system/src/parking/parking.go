package parking

import (
	"log"
	"ticket"
	"vehicle"
)

const (
	TWO_WHEELER  = 2
	FOUR_WHEELER = 4
)

type Capacity struct {
	Total_four_wheeler_space int
	Total_two_wheeler_space  int
}
type Parking struct {
	Capacity     *Capacity                `json:"capacity"`
	Parking_info map[int]*vehicle.Vehicle `json:"parking_info"`
}

func (p *Parking) is_parking_available(vehicle_type int) bool {
	switch vehicle_type {

	case TWO_WHEELER:
		return p.Capacity.Total_two_wheeler_space != 0

	case FOUR_WHEELER:
		return p.Capacity.Total_four_wheeler_space != 0

	default:
		log.Fatal("Vehile type: ", vehicle_type, " is not allowed")
		return false
	}
}

func (p *Parking) reserve_parking_space(vehicle_type int) int {
	var space_number int

	switch vehicle_type {

	case TWO_WHEELER:
		space_number = p.Capacity.Total_two_wheeler_space
		p.Capacity.Total_two_wheeler_space -= 1

	case FOUR_WHEELER:
		space_number = p.Capacity.Total_four_wheeler_space
		p.Capacity.Total_four_wheeler_space -= 1

	default:
		log.Fatal("Vehile type: ", vehicle_type, " is not allowed")
	}
	return space_number
}

func (p *Parking) get_price(vehicle_type int) int {
	switch vehicle_type {

	case TWO_WHEELER:
		return 30

	case FOUR_WHEELER:
		return 50

	default:
		return 0
	}
}
func (p *Parking) Get_reserve_parking_space(number string, vehicle_type int) *vehicle.Vehicle {
	is_parking_available := p.is_parking_available(vehicle_type)

	if is_parking_available == false {
		log.Fatal("There is no parking space left of type: ", vehicle_type)
		return &vehicle.Vehicle{}
	}

	// Reserve the parking space
	space_number := p.reserve_parking_space(vehicle_type)
	price := p.get_price(vehicle_type)

	ticket_obj := &ticket.Ticket{
		Vehicle_no:    number,
		Price:         price,
		Parking_space: space_number,
	}

	vehile_obj := &vehicle.Vehicle{
		Vehicle_number: number,
		Ticket:         ticket_obj,
		Type:           vehicle_type,
	}

	return vehile_obj
}
