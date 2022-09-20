package dinning_hall_elem

import (
	"fmt"
	"math/rand"
	"time"
)

type tableState int

// states of a table
const (
	Free tableState = iota
	WaitToOrder
	WaitToServe
)

type Table struct {
	Id          int
	State       tableState
	ClientOrder Order
	//used to control when the table is free
	TableChan chan int
}

func (t *Table) GenerateOrdersForever() {

	for _ = range t.TableChan {
		//The table thinks before generating a order
		time.Sleep(TimeUnit * time.Duration(rand.Intn(50)+3))
		t.makeOrder()

		fmt.Printf("Table %d generated order: %+v \n", t.Id, t.ClientOrder)
	}

}

// function for table to make order
func (t *Table) makeOrder() {
	//t.lock.Lock() //unlock will be called when the order is served

	t.State = WaitToOrder
	t.ClientOrder = newOrder()
	//sending the order to waiters; a waiter which is free will take it
	OrdersChannel <- t.Id
}

var Tables = []Table{
	{Id: 1, State: Free, TableChan: make(chan int, 1)}, {Id: 2, State: Free, TableChan: make(chan int, 1)},
	{Id: 3, State: Free, TableChan: make(chan int, 1)}, {Id: 4, State: Free, TableChan: make(chan int, 1)},
	{Id: 5, State: Free, TableChan: make(chan int, 1)}, {Id: 6, State: Free, TableChan: make(chan int, 1)},
	{Id: 7, State: Free, TableChan: make(chan int, 1)}, {Id: 8, State: Free, TableChan: make(chan int, 1)},
	{Id: 9, State: Free, TableChan: make(chan int, 1)}, {Id: 10, State: Free, TableChan: make(chan int, 1)},
}
