package vehicle

import (
	"ticket"
)

type Vehicle struct {
	Vehicle_number string `json:"vehicle_no"`
	Type           int    `json:"type"`
	Ticket         *ticket.Ticket
}
