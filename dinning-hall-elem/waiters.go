package dinning_hall_elem

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
	Lock        *sync.Mutex
	ordersChan  chan int
}

type kitchenFoodInf struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

func (w waiter) sayPhrase() {
	fmt.Printf("%s", w.catchPhrase)
}

// aici sa fie select cu channeles si true for loop
func (w *waiter) LookUpOrders() {

	//lock trebuie pentru a implementa ca chelnerul indata sa livreze comanda
	// chelnerii concomitent primesc din orders canal

	for i := range OrdersChannel {
		w.Lock.Lock()
		//DEFER DOES NOT WORK HERE BECAUSE THE FUNCTION NEVER RETURNS

		w.takeOrder(i)
		w.sendOrder()
		//w.free = true
		w.Lock.Unlock()
	}

}

//function to take table's order and prepare it for sending to kitchen
func (w *waiter) takeOrder(tableId int) {
	var table = Tables[tableId-1]
	var ord = sentOrd{
		OrderId:    table.ClientOrder.Id,
		TableId:    tableId,
		WaiterId:   w.id,
		Items:      table.ClientOrder.Items,
		Priority:   table.ClientOrder.Priority,
		MaxWait:    table.ClientOrder.MaxWait,
		PickUpTime: time.Now(),
	}
	w.takenOrder = ord
}

//function to send the order to kitchen using Post request
func (w *waiter) sendOrder() {
	reqBody, err := json.Marshal(w.takenOrder)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	resp, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewBuffer(reqBody))
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

	Tables[w.takenOrder.TableId-1].State = WaitToServe

	//for test, immediate serving

	//tables[w.takenOrder.TableId-1].state = free
	//tables[w.takenOrder.TableId-1].clientOrder = order{}
	//tables[w.takenOrder.TableId-1].lock.Unlock()
}

var Waiters = []waiter{
	{id: 1, catchPhrase: "Hii, i am Mikee", free: true, Lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 2, catchPhrase: "Finally got some tips from t a table!", free: true, Lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 3, catchPhrase: "Oh, what a grumpy client", free: true, Lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
	{id: 4, catchPhrase: "I love my work!", free: true, Lock: &sync.Mutex{}, ordersChan: make(chan int, 10)},
}
