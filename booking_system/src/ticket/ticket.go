package ticket

import (
	"time"
)

type Ticket struct {
	Vehicle_no    string
	Entry_time    time.Time
	Price         int
	Parking_space int
}
