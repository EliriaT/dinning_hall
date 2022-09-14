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

	//waiters[cookedOrder.WaiterId-1].ordersChan <- cookedOrder.TableId

	// i can use the waiter using lock, but it is not fastest way because it does not ensure that
	//as soon as possible, immediately, the waiter will serve
	//Lock is used to ensure that as soon as the waiter gets free,, it will send serve the order
	dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].Lock.Lock()
	log.Printf("Order with ID %d, was served by waiter %d. Details: %+v", cookedOrder.OrderId, cookedOrder.WaiterId, cookedOrder)
	dinning_hall_elem.Tables[cookedOrder.TableId-1].State = dinning_hall_elem.Free
	dinning_hall_elem.Tables[cookedOrder.TableId-1].ClientOrder = dinning_hall_elem.Order{}
	dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].Lock.Unlock()

	c.IndentedJSON(http.StatusCreated, cookedOrder)
	//tables[cookedOrder.TableId].lock.Unlock()
}

func main() {

	rand.Seed(time.Now().UnixNano())

	dinning_hall_elem.AiOrder.SetId(1)
	dinning_hall_elem.Init()
	//wg.Wait()
	router := gin.Default()
	router.GET("/foods", getFoods)
	router.POST("/distribution", serveOrder)
	for i, _ := range dinning_hall_elem.Waiters {
		//i should check if the waiter is free
		go dinning_hall_elem.Waiters[i].LookUpOrders()
	}

	router.Run(":8082")

	//fmt.Printf("hi\n")
}
