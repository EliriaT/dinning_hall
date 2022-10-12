package dinning_hall_elem

import (
	"sync"
	"time"
)

const (
	TimeUnit = time.Duration(float64(time.Millisecond) * 25)
	//OrdersLimit = 1000
	MaxFoods = 10
	minFoods = 1
)

type tableState int

// states of a table
const (
	Free tableState = iota
	WaitToOrder
	WaitToServe
)

var (
	Port           string
	DinningHallUrl string
	KitchenURL     string
	//KitchenURL = "http://localhost:8080/"
	ManagerURL     string
	RestaurantId   int
	RestaurantName string
	nrTables       int
	NrFoods        int
	Rating         float32

	Tables  []Table
	Waiters []waiter
	Foods   []Food

	OrdersChannel chan int
	// Used to concurrently safe append to the RaitingSlice
	MarkMutex  = &sync.Mutex{}
	markLength = 0
	sum        = 0

	//ListOfOnlineCookedOrders
	OnlineCookedOrder = make([]ReceivedOrd, 0, 100)
)
