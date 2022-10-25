package versionTwoElems

import (
	dinning_hall_elem "github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"time"
)

type OnlineReceivedOrder struct {
	Id          int       `json:"id,omitempty"`
	Items       []int     `json:"items"`
	Priority    int       `json:"priority"`
	MaxWait     float32   `json:"max_wait"`
	CreatedTime time.Time `json:"created_time"`
}

type OnlineResponseOrder struct {
	RestaurantId         int       `json:"restaurant_id"`
	OrderId              int       `json:"order_id"`
	EstimatedWaitingTime int       `json:"estimated_waiting_time"`
	CreatedTime          time.Time `json:"created_time"`
	RegisteredTime       time.Time `json:"registered_time"`
}

type OnlineCookedOrder struct {
	OrderId        int                                `json:"order_id"`
	IsReady        bool                               `json:"is_ready"`
	EstimatedTime  int                                `json:"estimated_waiting_time"`
	Priority       int                                `json:"priority"`
	MaxWait        float32                            `json:"max_wait"`
	CreatedTime    time.Time                          `json:"created_time"`
	RegisteredTime time.Time                          `json:"registered_time"`
	PreparedTime   time.Time                          `json:"prepared_time"`
	CookingTime    time.Duration                      `json:"cooking_time"`
	CookingDetails []dinning_hall_elem.KitchenFoodInf `json:"cooking_details"`
}

// ListOfOnlineCookedOrders where the cooked orders are saved after receiving from kitchen
var OnlineCookedOrdersMap = make(map[int]OnlineCookedOrder)
