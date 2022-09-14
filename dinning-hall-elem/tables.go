package dinning_hall_elem

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type tableState int

// states of a table
const (
	Free tableState = iota
	WaitToOrder
	WaitToServe
)
const TimeUnit = 50 * time.Millisecond

// mutex to control when the table can generate order?
type Table struct {
	Id          int
	State       tableState
	ClientOrder Order
	//Probably will need in future?
	Lock *sync.Mutex
}

// function for table to make order
func (t *Table) makeOrder() {
	//t.lock.Lock() //unlock will be called when the order is served
	time.Sleep(TimeUnit * time.Duration(rand.Intn(2)+2))
	t.State = WaitToOrder
	t.ClientOrder = newOrder()
	OrdersChannel <- t.Id //?
}

var Tables = []Table{
	{Id: 1, State: Free, Lock: &sync.Mutex{}}, {Id: 2, State: Free, Lock: &sync.Mutex{}},
	{Id: 3, State: Free, Lock: &sync.Mutex{}}, {Id: 4, State: Free, Lock: &sync.Mutex{}},
	{Id: 5, State: Free, Lock: &sync.Mutex{}}, {Id: 6, State: Free, Lock: &sync.Mutex{}},
	{Id: 7, State: Free, Lock: &sync.Mutex{}}, {Id: 8, State: Free, Lock: &sync.Mutex{}},
	{Id: 9, State: Free, Lock: &sync.Mutex{}}, {Id: 10, State: Free, Lock: &sync.Mutex{}},
}

// init function to initialize first orders
func Init() {
	//defer wg.Done() IT CREATES A COPY OF TABLE IN T
	for i, _ := range Tables {
		if rand.Intn(2) == 1 {
			Tables[i].makeOrder()
			fmt.Printf("%+v of table %d \n", Tables[i].ClientOrder, Tables[i].Id)
		}
	}
	log.Printf("Init finished")
}
