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
