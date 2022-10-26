package versionTwoElems

type OrderRating struct {
	OrderId       int `json:"order_id"`
	Rating        int `json:"rating"`
	EstimatedTime int `json:"estimated_waiting_time"`
	WaitedTime    int `json:"waiting_time"`
}

type RestaurantRating struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float32 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
