package dinning_hall_elem

import "time"

const (
	TimeUnit    = time.Duration(float64(time.Millisecond) * 10)
	OrdersLimit = 50
	maxFoods    = 10
	minFoods    = 1
	nrFoods     = 13
	URL         = "http://localhost:8080/order"
)

type tableState int

// states of a table
const (
	Free tableState = iota
	WaitToOrder
	WaitToServe
)
