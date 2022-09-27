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

func (w *waiter) serveOrder(cookedOrder ReceivedOrd) {

	serveTime := time.Since(cookedOrder.PickUpTime)
	//
	orderRaiting := giveOrderStars(serveTime, cookedOrder.MaxWait)

	log.Printf("Order with ID %d, was SERVED by waiter %d. Details: %+v \n ", cookedOrder.OrderId, cookedOrder.WaiterId, cookedOrder)

	log.Printf("RAITING OF ORDER IS %d, MAX TIME IS %f, SERVE TIME IS %v", orderRaiting, cookedOrder.MaxWait, serveTime)
	Tables[cookedOrder.TableId-1].State = Free
	Tables[cookedOrder.TableId-1].ClientOrder = Order{}

	//freeing the table
	MarkMutex.Lock()
	OrderMarks = append(OrderMarks, orderRaiting)
	MarkMutex.Unlock()

	//For allowing to limit work time of the restaurant
	if len(OrderMarks) == OrdersLimit {
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
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(reqBody))
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
