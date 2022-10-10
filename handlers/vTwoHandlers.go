package handlers

import (
	"encoding/json"
	dinning_hall_elem "github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/EliriaT/dinning_hall/versionTwoElems"
	"net/http"
	"time"
)

func ReceiveOnlineOrder(w http.ResponseWriter, r *http.Request) {
	var receivedOrder versionTwoElems.OnlineReceivedOrder
	var responseOrder versionTwoElems.OnlineResponseOrder
	var kitchenInform versionTwoElems.KitchenInfo

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

	responseOrder.RestaurantId = dinning_hall_elem.RestaurantId
	responseOrder.OrderId = receivedOrder.Id
	responseOrder.EstimatedWaitingTime = estimatedTime
	responseOrder.CreatedTime = receivedOrder.CreatedTime
	responseOrder.RegisteredTime = time.Now()

	defer r.Body.Close()
	jsonCookedOrder, _ := json.Marshal(responseOrder)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
	//FINISH
}
