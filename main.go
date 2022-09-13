package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func getFoods(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, foods)
}

// the waiter should be free, it should be added at a list of needed to serve order, or to notify the waiter, to implement free waiter with mutex
func serveOrder(c *gin.Context) {
	var cookedOrder receivedOrd
	if err := c.BindJSON(&cookedOrder); err != nil {
		log.Printf(err.Error())
	}
	tables[cookedOrder.TableId].state = free
	tables[cookedOrder.TableId].clientOrder = order{}
	tables[cookedOrder.TableId].lock.Unlock()

	log.Printf("Order with ID %d, was served", cookedOrder.OrderId)
	c.IndentedJSON(http.StatusCreated, cookedOrder)
}

func main() {
	//wg := new(sync.WaitGroup)
	//wg.Add(1)
	rand.Seed(time.Now().UnixNano())
	Init()
	//wg.Wait()
	router := gin.Default()
	router.GET("/foods", getFoods)
	router.POST("/distribution", serveOrder)
	for i, _ := range waiters {
		//i should check if the waiter is free
		go waiters[i].lookUpOrders()
	}

	router.Run("localhost:8080")

	//fmt.Printf("hi\n")
}
