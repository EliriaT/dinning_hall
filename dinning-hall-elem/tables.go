package dinning_hall_elem

import (
	"fmt"
	"log"
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

// mutex to control when the table can generate order?
type Table struct {
	Id          int
	State       tableState
	ClientOrder Order
	//this channel has data when the table becomes free, otherwise is in a wait to generate order
	TableChan chan int
	//Probably will need in future?
	//Lock *sync.Mutex
}

// i can use the free state atribute, but then i will have to use a lock, because in the same time, other go routine may access table's state; better to use channels
func (t *Table) GenerateOrdersForever() {

	for _ = range t.TableChan {
		//The table thinks before generating a order
		time.Sleep(TimeUnit * time.Duration(rand.Intn(150)+10))
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

// init function to initialize first orders
func Init() {

	for i, _ := range Tables {
		// this means table is free
		Tables[i].TableChan <- 1
	}

	//random number of tables up to 5 can at start generate order
	nrTablesInit := rand.Intn(5) + 1

	//Get random ID's of tables shuffled
	randTableInit := rand.Perm(10)

	//Get only the first n random ID's
	randTableInit = randTableInit[0:nrTablesInit]

	for _, i := range randTableInit {
		Tables[i].makeOrder()
		<-Tables[i].TableChan
		fmt.Printf("Table %d generated order: %+v \n", i+1, Tables[i].ClientOrder)

	}
	log.Printf("Init finished")
}

var Tables = []Table{
	{Id: 1, State: Free, TableChan: make(chan int, 1)}, {Id: 2, State: Free, TableChan: make(chan int, 1)},
	{Id: 3, State: Free, TableChan: make(chan int, 1)}, {Id: 4, State: Free, TableChan: make(chan int, 1)},
	{Id: 5, State: Free, TableChan: make(chan int, 1)}, {Id: 6, State: Free, TableChan: make(chan int, 1)},
	{Id: 7, State: Free, TableChan: make(chan int, 1)}, {Id: 8, State: Free, TableChan: make(chan int, 1)},
	{Id: 9, State: Free, TableChan: make(chan int, 1)}, {Id: 10, State: Free, TableChan: make(chan int, 1)},
}
