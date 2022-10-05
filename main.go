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

func main() {

	runtime.GOMAXPROCS(1)
	rand.Seed(time.Now().UnixNano())
	//

	dinning_hall_elem.AiOrder.SetId(1)
	dinning_hall_elem.Init()
	log.Println("finished")
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

}
