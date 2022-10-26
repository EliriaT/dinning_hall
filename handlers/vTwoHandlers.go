package handlers

import (
	"encoding/json"
	dinning_hall_elem "github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/EliriaT/dinning_hall/versionTwoElems"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ReceiveOnlineOrder(w http.ResponseWriter, r *http.Request) {
	var receivedOrder versionTwoElems.OnlineReceivedOrder
	var responseOrder versionTwoElems.OnlineResponseOrder
	var kitchenInform versionTwoElems.KitchenInfo
	var registeredOrder versionTwoElems.OnlineCookedOrder

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receivedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receivedOrder.Id = dinning_hall_elem.AiOrder.ID()

	//sending the order immediately to the kitchen
	kitchenInform = versionTwoElems.SendOnlineOrder(receivedOrder)

	//calculating estimating time
	estimatedTime := versionTwoElems.CalculateEstimatedTime(receivedOrder, kitchenInform)
	//if float32(estimatedTime) > receivedOrder.MaxWait {
	//	estimatedTime = int(receivedOrder.MaxWait) - 2
	//}

	//TODO Orderul inregistrat trebuie salvat intrun registru
	responseOrder.RestaurantId = dinning_hall_elem.RestaurantId
	responseOrder.OrderId = receivedOrder.Id
	responseOrder.EstimatedWaitingTime = estimatedTime
	responseOrder.CreatedTime = receivedOrder.CreatedTime
	responseOrder.RegisteredTime = time.Now()

	defer r.Body.Close()
	jsonCookedOrder, _ := json.Marshal(responseOrder)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)

	registeredOrder.OrderId = receivedOrder.Id
	registeredOrder.IsReady = false
	registeredOrder.EstimatedTime = estimatedTime
	registeredOrder.Priority = receivedOrder.Priority
	registeredOrder.MaxWait = receivedOrder.MaxWait
	registeredOrder.CreatedTime = receivedOrder.CreatedTime
	registeredOrder.RegisteredTime = responseOrder.RegisteredTime
	registeredOrder.PreparedTime = time.Now()
	registeredOrder.CookingTime = 0
	registeredOrder.CookingDetails = nil

	versionTwoElems.OnlineCookedOrdersMap[registeredOrder.OrderId] = registeredOrder

}

func SendOnlineOrderToClient(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal("Could not convert id to int")
	}

	cookedOnlineOrder, ok := versionTwoElems.OnlineCookedOrdersMap[orderId]
	if ok == false {
		log.Printf("---- Order %d is not found in the register. Some strange error----", orderId)
	}

	if cookedOnlineOrder.IsReady == true {
		log.Printf("---- Order %d is ready.Pick up can be done----", orderId)
		//return
	}

	jsonCookedOrder, _ := json.Marshal(cookedOnlineOrder)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)

}

func ReceiveClientRating(w http.ResponseWriter, r *http.Request) {
	var receivedClientRating versionTwoElems.OrderRating
	var sendAvgRating versionTwoElems.RestaurantRating

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receivedClientRating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dinning_hall_elem.MarkMutex.Lock()
	dinning_hall_elem.MarkLength++
	avg := dinning_hall_elem.CalculateAverage(receivedClientRating.Rating)
	log.Println("-----Average is : ", avg, "-----")
	dinning_hall_elem.MarkMutex.Unlock()

	sendAvgRating.RestaurantId = dinning_hall_elem.RestaurantId
	sendAvgRating.RestaurantAvgRating = float32(avg)
	sendAvgRating.PreparedOrders = dinning_hall_elem.AiOrder.Id

	jsonSendAvgRating, _ := json.Marshal(sendAvgRating)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonSendAvgRating)
}
