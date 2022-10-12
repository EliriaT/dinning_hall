package dinning_hall_elem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func initTable() {
	Tables = make([]Table, nrTables, nrTables)

	for i := 0; i < nrTables; i++ {
		Tables[i].Id = i + 1
		Tables[i].State = Free
		Tables[i].TableChan = make(chan int, 1)
	}
}

func initWaiter() {
	file, err := os.Open("./jsonConfig/waiters.json")
	if err != nil {
		log.Fatal("Error opening waiters.json", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	_ = json.Unmarshal(byteValue, &Waiters)

	for i := range Waiters {
		Waiters[i].CookedOrdersChan = make(chan ReceivedOrd, 10)
	}
	//log.Println("Waiters: ", Waiters)
}

func initiate_Foods() {
	file, err := os.Open("./jsonConfig/foods.json")
	if err != nil {
		log.Fatal("Error opening foods.json", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	_ = json.Unmarshal(byteValue, &Foods)

	NrFoods = len(Foods)

}

func initiate_Congif() {
	var config map[string]string

	file, err := os.Open("./jsonConfig/config.json")
	if err != nil {
		log.Fatal("Error opening config.json ", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	_ = json.Unmarshal(byteValue, &config)

	DinningHallUrl = config["my_address"]
	KitchenURL = config["kitchen_address"]
	ManagerURL = config["food_manager_address"]
	Port = config["listenning_port"]
	RestaurantName = config["restaurant_name"]
	RestaurantId, err = strconv.Atoi(config["resaurant_id"])
	nrTables, err = strconv.Atoi(config["nr_tables"])
	OrdersChannel = make(chan int, nrTables)

	if err != nil {
		log.Fatal("Error int conversion")
	}

}

// init function to initialize first orders
func Init() {
	initiate_Congif()

	initTable()

	initWaiter()

	initiate_Foods()

	for i := range Tables {
		// this means table is free
		Tables[i].TableChan <- 1
	}

	//random number of tables up to 5 can at start generate order
	nrTablesInit := rand.Intn(3) + 1

	//Get random ID's of tables shuffled
	randTableInit := rand.Perm(nrTables)

	//Get only the first n random ID's
	randTableInit = randTableInit[0:nrTablesInit]

	for _, i := range randTableInit {
		Tables[i].makeOrder()
		<-Tables[i].TableChan
		fmt.Printf("Table %d generated order: %+v \n", i+1, Tables[i].ClientOrder)

	}

	log.Printf("Init finished")
}
