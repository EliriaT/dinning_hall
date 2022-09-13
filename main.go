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
		return
	}

	//waiters[cookedOrder.WaiterId-1].ordersChan <- cookedOrder.TableId

	// i can use the waiter using lock, but it is not fastest way because it does not ensure that
	//as soon as possible, immediately, the waiter will serve

	waiters[cookedOrder.WaiterId-1].lock.Lock()
	log.Printf("Order with ID %d, was served by waiter %d. Details: %+v", cookedOrder.OrderId, cookedOrder.WaiterId, cookedOrder)
	tables[cookedOrder.TableId-1].state = free
	tables[cookedOrder.TableId-1].clientOrder = order{}
	waiters[cookedOrder.WaiterId-1].lock.Unlock()

	c.IndentedJSON(http.StatusCreated, cookedOrder)
	//tables[cookedOrder.TableId].lock.Unlock()
}

func main() {
	//wg := new(sync.WaitGroup)
	//wg.Add(1)
	rand.Seed(time.Now().UnixNano())
	aiOrder.SetId(1)
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
