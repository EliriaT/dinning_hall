package main

import (
	"encoding/json"
	"github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/gorilla/mux"
	"runtime"

	//_ "go.uber.org/automaxprocs"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func getFoods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonFoods, err := json.Marshal(dinning_hall_elem.Foods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//by default sends 200
	w.Write(jsonFoods)
}

// the waiter should be free, it should be added at a list of needed to serve order, or to notify the waiter, to implement free waiter with mutex
func serveOrder(w http.ResponseWriter, r *http.Request) {
	var cookedOrder dinning_hall_elem.ReceivedOrd
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cookedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].CookedOrdersChan <- cookedOrder

	defer r.Body.Close()
	jsonCookedOrder, _ := json.Marshal(cookedOrder)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
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
	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UnixNano())
	//
	dinning_hall_elem.AiOrder.SetId(1)
	dinning_hall_elem.Init()

	r := mux.NewRouter()
	r.HandleFunc("/", getFoods).Methods("GET")
	r.HandleFunc("/distribution", serveOrder).Methods("POST")
	for i, _ := range dinning_hall_elem.Waiters {

		go dinning_hall_elem.Waiters[i].Work()
	}

	for i, _ := range dinning_hall_elem.Tables {
		go dinning_hall_elem.Tables[i].GenerateOrdersForever()
	}

	log.Println("Dinning Hall server started..")
	log.Fatal(http.ListenAndServe(":8082", r))

	//fmt.Printf("hi\n")
}
