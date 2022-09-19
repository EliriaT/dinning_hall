package main

import (
	"github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func getFoods(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dinning_hall_elem.Foods)
}

// the waiter should be free, it should be added at a list of needed to serve order, or to notify the waiter, to implement free waiter with mutex
func serveOrder(c *gin.Context) {
	var cookedOrder dinning_hall_elem.ReceivedOrd
	if err := c.BindJSON(&cookedOrder); err != nil {
		log.Printf(err.Error())
		return
	}

	dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].CookedOrdersChan <- cookedOrder

	// i can use the waiter using lock, but it is not fastest way because it does not ensure that
	//as soon as possible, immediately, the waiter will serve, here I want to use select with default , then select in the default

	//Lock is used to ensure that as soon as the waiter gets free, it will send serve the order

	//dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].Lock.Lock()
	//serveTime := time.Since(cookedOrder.PickUpTime)
	//
	//
	//orderRaiting := giveOrderStars(serveTime, cookedOrder.MaxWait)
	//
	//log.Printf("Order with ID %d, was SERVED by waiter %d. Details: %+v \n ", cookedOrder.OrderId, cookedOrder.WaiterId, cookedOrder)
	//
	//log.Printf("RAITING OF ORDER IS %d, MAX TIME IS %f, SERVE TIME IS %v", orderRaiting, cookedOrder.MaxWait, serveTime)
	//dinning_hall_elem.Tables[cookedOrder.TableId-1].State = dinning_hall_elem.Free
	//dinning_hall_elem.Tables[cookedOrder.TableId-1].ClientOrder = dinning_hall_elem.Order{}
	//
	//dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].Lock.Unlock()
	//dinning_hall_elem.Tables[cookedOrder.TableId-1].TableChan <- 1

	c.IndentedJSON(http.StatusCreated, cookedOrder)

}

func giveOrderStars(serveTime time.Duration, maxWait float64) int {
	//serveTimeMillisec := float64(serveTime)*1000 //time in milliseconds
	serveTimeNonUnit := float64(serveTime) / float64(dinning_hall_elem.TimeUnit)
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

func main() {

	rand.Seed(time.Now().UnixNano())
	//
	dinning_hall_elem.AiOrder.SetId(1)
	dinning_hall_elem.Init()

	router := gin.Default()
	router.GET("/foods", getFoods)
	router.POST("/distribution", serveOrder)
	for i, _ := range dinning_hall_elem.Waiters {

		go dinning_hall_elem.Waiters[i].Work()
	}

	for i, _ := range dinning_hall_elem.Tables {
		go dinning_hall_elem.Tables[i].GenerateOrdersForever()
	}

	router.Run(":8082")

	//fmt.Printf("hi\n")
}
