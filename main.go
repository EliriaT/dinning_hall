package main

import (
	"github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/EliriaT/dinning_hall/handlers"
	"github.com/EliriaT/dinning_hall/versionTwoElems"
	"github.com/gorilla/mux"
	"runtime"

	//_ "go.uber.org/automaxprocs"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	runtime.GOMAXPROCS(1)
	rand.Seed(time.Now().UnixNano())
	//TODO REGISTER RESTAURANT

	dinning_hall_elem.AiOrder.SetId(1)
	dinning_hall_elem.Init()
	versionTwoElems.RegisterRestaurant()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.GetFoods).Methods("GET")
	r.HandleFunc("/distribution", handlers.ServeOrder).Methods("POST")

	s := r.PathPrefix("/v2").Subrouter()
	s.HandleFunc("/order", handlers.ReceiveOnlineOrder)

	//s.HandleFunc("/order/{id:[0-9]+}", handlers.ReceiveOnlineOrder)

	for i, _ := range dinning_hall_elem.Waiters {

		go dinning_hall_elem.Waiters[i].Work()
	}

	for i, _ := range dinning_hall_elem.Tables {
		go dinning_hall_elem.Tables[i].GenerateOrdersForever()
	}

	log.Println("Dinning Hall server started..")
	log.Fatal(http.ListenAndServe(":8082", r))

}
