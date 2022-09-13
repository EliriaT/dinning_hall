package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type waiter struct {
	id          int
	catchPhrase string
	takenOrder  sentOrd
	free        bool
	lock        *sync.Mutex
	ordersChan  chan int
}

type receivedOrd struct {
	OrderId        int              `json:"order_id"`
	TableId        int              `json:"table_id"`
	WaiterId       int              `json:"waiter_id"`
	Items          []int            `json:"items"`
	Priority       int              `json:"priority"`
	MaxWait        int              `json:"max_wait"`
	PickUpTime     time.Time        `json:"pick_up_time"`
	CookingTime    time.Duration    `json:"cooking_time"`
	CookingDetails []kitchenFoodInf `json:"cooking_details"`
}
type kitchenFoodInf struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}
type sentOrd struct {
	OrderId    int       `json:"order_id"`
	TableId    int       `json:"table_id"`
	WaiterId   int       `json:"waiter_id"`
	Items      []int     `json:"items"`
	Priority   int       `json:"priority"`
	MaxWait    int       `json:"max_wait"`
	PickUpTime time.Time `json:"pick_up_time"`
}

func (w waiter) sayPhrase() {
	fmt.Printf("%s", w.catchPhrase)
}

// aici sa fie select cu channeles si true for loop
func (w *waiter) lookUpOrders() {

	//lock trebuie pentru a implementa ca chelnerul indata sa livreze comanda
	//dar el chelnerul nu poate intra in for loop daca el nu e free, sau poate ? poate
	// chelnerii concomitent primesc din orders canal

	for i := range ordersChannel {
		w.lock.Lock()
		//DEFER DOES NOT WORK BECAUSE THE FUNCTION NEVER RETURNS
		//defer w.lock.Unlock() //after it sent the order, it can continue receiving other orders from tables
		//w.free = false
		w.takeOrder(i)
		w.sendOrder()
		//w.free = true
		w.lock.Unlock()
	}

}

func (w *waiter) takeOrder(tableId int) {
	var table = tables[tableId-1]
	var ord = sentOrd{
		OrderId:    table.clientOrder.Id,
		TableId:    tableId,
		WaiterId:   w.id,
		Items:      table.clientOrder.Items,
		Priority:   table.clientOrder.Priority,
		MaxWait:    table.clientOrder.MaxWait,
		PickUpTime: time.Now(),
	}
	w.takenOrder = ord
}

func (w *waiter) sendOrder() {
	reqBody, err := json.Marshal(w.takenOrder)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	resp, err := http.Post("http://localhost:8082/order", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Request Failed: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body) // Log the request body
	//if err != nil {
	//	log.Printf("Can't read the response body %s", err.Error())
	//	return
	//}
	//bodyString := string(body)
	//log.Print(bodyString)
	log.Printf("The order with id %d was sent to Kitchen by waiter %d. Details: %+v", w.takenOrder.OrderId, w.takenOrder.WaiterId, w.takenOrder) // Unmarshal result
	//aici sa modific statusul la waiting to serve la table
	tables[w.takenOrder.TableId-1].state = waitToServe

	//for test, immediare serving
	//log.Printf("Order with ID %d, was served by waiter %d. Details: %v", w.takenOrder.OrderId, w.takenOrder.WaiterId, w.takenOrder)
	//tables[w.takenOrder.TableId-1].state = free
	//tables[w.takenOrder.TableId-1].clientOrder = order{}
	//tables[w.takenOrder.TableId-1].lock.Unlock()
}

var waiters = []waiter{
	{id: 1, catchPhrase: "Hii, i am Mikee", free: true, lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 2, catchPhrase: "Finally got some tips from t a table!", free: true, lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 3, catchPhrase: "Oh, what a grumpy client", free: true, lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 4, catchPhrase: "I love my work!", free: true, lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
}
