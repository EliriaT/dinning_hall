package dinning_hall_elem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type waiter struct {
	id               int
	catchPhrase      string
	takenOrder       sentOrd
	free             bool
	Lock             *sync.Mutex
	CookedOrdersChan chan ReceivedOrd
}

var MarkMutex = &sync.Mutex{}

type kitchenFoodInf struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

func (w waiter) sayPhrase() {
	fmt.Printf("%s", w.catchPhrase)
}

// aici sa fie select cu channeles si true for loop
//func (w *waiter) LookUpOrders() {
//
//	//lock trebuie pentru a implementa ca chelnerul indata sa livreze comanda
//	// chelnerii concomitent primesc din orders canal
//
//	for i := range OrdersChannel {
//		w.Lock.Lock()
//
//		w.takeOrder(i)
//		w.sendOrder()
//		//w.free = true
//		w.Lock.Unlock()
//	}
//
//}

func (w *waiter) Work() {
	for {
		select {
		case kitchOrder := <-w.CookedOrdersChan:
			w.serveOrder(kitchOrder)
		default:

			select {
			case kitchOrder := <-w.CookedOrdersChan:
				w.serveOrder(kitchOrder)
			case hallOrder := <-OrdersChannel:

				select {
				case kitchOrder := <-w.CookedOrdersChan:
					w.serveOrder(kitchOrder)

					//w.Lock.Lock()

					w.takeOrder(hallOrder)
					w.sendOrder()

					//w.Lock.Unlock()
				default:

					//w.Lock.Lock()

					w.takeOrder(hallOrder)
					w.sendOrder()

					//w.Lock.Unlock()
				}
			}
		}
	}
}

func calculateAverage() float64 {
	sum := 0
	for _, mark := range OrderMarks {
		sum += mark
	}
	avg := float64(sum) / float64(len(OrderMarks))
	return avg
}

func (w *waiter) serveOrder(cookedOrder ReceivedOrd) {
	//w.Lock.Lock()
	serveTime := time.Since(cookedOrder.PickUpTime)
	//
	orderRaiting := giveOrderStars(serveTime, cookedOrder.MaxWait)

	log.Printf("Order with ID %d, was SERVED by waiter %d. Details: %+v \n ", cookedOrder.OrderId, cookedOrder.WaiterId, cookedOrder)

	log.Printf("RAITING OF ORDER IS %d, MAX TIME IS %f, SERVE TIME IS %v", orderRaiting, cookedOrder.MaxWait, serveTime)
	Tables[cookedOrder.TableId-1].State = Free
	Tables[cookedOrder.TableId-1].ClientOrder = Order{}

	//w.Lock.Unlock()
	//freeing the table
	MarkMutex.Lock()
	OrderMarks = append(OrderMarks, orderRaiting)
	MarkMutex.Unlock()
	if len(OrderMarks) == 15 {
		avg := calculateAverage()
		log.Println("Program Terminating; Restaurant's average is : ", avg)
		os.Exit(0)
	}

	Tables[cookedOrder.TableId-1].TableChan <- 1

}

// function to take table's order and prepare it for sending to kitchen
func (w *waiter) takeOrder(tableId int) {
	//taking of order takes time
	time.Sleep(TimeUnit * time.Duration(rand.Intn(8)+5))

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

// function to send the order to kitchen using Post request
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

	log.Printf("The order with id %d was SENT to Kitchen by waiter %d. Details: %+v\n", w.takenOrder.OrderId, w.takenOrder.WaiterId, w.takenOrder) // Unmarshal result

	Tables[w.takenOrder.TableId-1].State = WaitToServe

}

var Waiters = []waiter{
	{id: 1, catchPhrase: "Hii, i am Mikee", free: true, Lock: &sync.Mutex{}, CookedOrdersChan: make(chan ReceivedOrd, 10)},
	{id: 2, catchPhrase: "Finally got some tips from t a table!", free: true, Lock: &sync.Mutex{}, CookedOrdersChan: make(chan ReceivedOrd, 10)},
	{id: 3, catchPhrase: "Oh, what a grumpy client", free: true, Lock: &sync.Mutex{}, CookedOrdersChan: make(chan ReceivedOrd, 10)},
	{id: 4, catchPhrase: "I love my work!", free: true, Lock: &sync.Mutex{}, CookedOrdersChan: make(chan ReceivedOrd, 10)},
}

func giveOrderStars(serveTime time.Duration, maxWait float64) int {
	//serveTimeMillisec := float64(serveTime)*1000 //time in milliseconds
	serveTimeNonUnit := float64(serveTime) / float64(TimeUnit)
	//int(serveTime)/int(TimeUnit)
	//maxWait=maxWait*
	switch {
	case serveTimeNonUnit < maxWait:
		return 5
	case serveTimeNonUnit < maxWait*1.1:
		return 4
	case serveTimeNonUnit < maxWait*1.2:
		return 3
	case serveTimeNonUnit < maxWait*1.3:
		return 2
	case serveTimeNonUnit < maxWait*1.4:
		return 1
	default:
		return 0
	}
}
