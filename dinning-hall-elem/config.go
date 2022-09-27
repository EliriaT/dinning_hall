package dinning_hall_elem

import "time"

const (
	TimeUnit    = time.Duration(float64(time.Millisecond) * 50)
	OrdersLimit = 100

	URL = "http://localhost:8080/order"
)
