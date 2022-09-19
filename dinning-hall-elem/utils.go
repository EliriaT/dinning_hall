package dinning_hall_elem

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func calculateAverage() float64 {
	sum := 0
	for _, mark := range OrderMarks {
		sum += mark
	}
	avg := float64(sum) / float64(len(OrderMarks))
	return avg
}

func giveOrderStars(serveTime time.Duration, maxWait float64) int {
	//serveTimeMillisec := float64(serveTime)*1000 //time in milliseconds
	serveTimeNonUnit := float64(serveTime) / float64(TimeUnit)

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

// init function to initialize first orders
func Init() {

	for i, _ := range Tables {
		// this means table is free
		Tables[i].TableChan <- 1
	}

	//random number of tables up to 5 can at start generate order
	nrTablesInit := rand.Intn(5) + 1

	//Get random ID's of tables shuffled
	randTableInit := rand.Perm(10)

	//Get only the first n random ID's
	randTableInit = randTableInit[0:nrTablesInit]

	for _, i := range randTableInit {
		Tables[i].makeOrder()
		<-Tables[i].TableChan
		fmt.Printf("Table %d generated order: %+v \n", i+1, Tables[i].ClientOrder)

	}
	log.Printf("Init finished")
}