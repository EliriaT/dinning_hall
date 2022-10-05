package dinning_hall_elem

import (
	"fmt"
	"math/rand"
	"time"
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
		time.Sleep(TimeUnit * time.Duration(rand.Intn(500)+30))
		t.makeOrder()

		fmt.Printf("Table %d generated order: %+v \n", t.Id, t.ClientOrder)
	}

}

// function for table to make order
func (t *Table) makeOrder() {

	t.State = WaitToOrder
	t.ClientOrder = newOrder()
	//sending the order to waiters; a waiter which is free will take it; All waiters are listening to the OrdersChannel
	OrdersChannel <- t.Id

}
