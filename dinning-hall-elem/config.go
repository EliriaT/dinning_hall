package dinning_hall_elem

import (
	"sync"
	"time"
)

const (
	TimeUnit    = time.Duration(float64(time.Millisecond) * 25)
	OrdersLimit = 1000
	MaxFoods    = 10
	minFoods    = 1
	NrFoods     = 13
	nrTables    = 10

	Rating = 5
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
	DinningHallUrl = "http://localhost:8082/"
	KitchenURL     = "http://kitchen1:8080/"
	//KitchenURL = "http://localhost:8080/"
	ManagerURL     = "http://localhost:8084/"
	RestaurantId   = 1
	RestaurantName = "La Placinte"

	Tables  []Table
	Waiters []waiter
	Foods   []Food
	//OrderMarks    = make([]int, 0, 100)
	OrdersChannel = make(chan int, nrTables)
	// Used to concurrently safe append to the RaitingSlice
	MarkMutex  = &sync.Mutex{}
	markLength = 0
	sum        = 0

	//ListOfOnlineCookedOrders
	OnlineCookedOrder = make([]ReceivedOrd, 0, 100)
)
