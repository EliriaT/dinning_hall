package handlers

import (
	"encoding/json"
	dinning_hall_elem "github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"github.com/EliriaT/dinning_hall/versionTwoElems"
	"net/http"
)

func GetFoods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonFoods, err := json.Marshal(dinning_hall_elem.Foods)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//by default sends 200
	w.Write(jsonFoods)
}

func ServeOrder(w http.ResponseWriter, r *http.Request) {
	var cookedOrder dinning_hall_elem.ReceivedOrd
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cookedOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//the order received from kitchen
	if cookedOrder.WaiterId == -1 {
		order := versionTwoElems.OnlineCookedOrdersMap[cookedOrder.OrderId]
		order.IsReady = true
		order.EstimatedTime = 0
		order.CookingTime = cookedOrder.CookingTime
		order.CookingDetails = cookedOrder.CookingDetails
		versionTwoElems.OnlineCookedOrdersMap[cookedOrder.OrderId] = order

		//dinning_hall_elem.OnlineCookedOrder = append(dinning_hall_elem.OnlineCookedOrder, cookedOrder)
	} else {
		dinning_hall_elem.Waiters[cookedOrder.WaiterId-1].CookedOrdersChan <- cookedOrder

	}

	defer r.Body.Close()
	//i think i do not need this as it can take unneccessary time
	jsonCookedOrder, _ := json.Marshal(cookedOrder)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCookedOrder)
}
