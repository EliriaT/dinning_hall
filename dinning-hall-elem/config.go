package dinning_hall_elem

import (
	"sync"
	"time"
)

const (
	TimeUnit    = time.Duration(float64(time.Millisecond) * 10)
	OrdersLimit = 100
	maxFoods    = 10
	minFoods    = 1
	nrFoods     = 13
	nrTables    = 10
	//URL         = "http://kitchen:8080/order"
	URL = "http://localhost:8080/order"
)

type tableState int

// states of a table
const (
	Free tableState = iota
	WaitToOrder
	WaitToServe
)

var (
	Tables        []Table
	Waiters       []waiter
	Foods         []food
	OrderMarks    = make([]int, 0, 100)
	OrdersChannel = make(chan int, nrTables)
	// Used to concurrently safe append to the RaitingSlice
	MarkMutex = &sync.Mutex{}
)
