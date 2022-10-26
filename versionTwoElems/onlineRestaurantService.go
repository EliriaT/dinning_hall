package versionTwoElems

import (
	"bytes"
	"encoding/json"
	"github.com/EliriaT/dinning_hall/dinning-hall-elem"
	"io"
	"log"
	"net/http"
)

func CalculateEstimatedTime(currentOrder OnlineReceivedOrder, kitchenInform KitchenInfo) int {
	var A, B, C, D, E, F, estimatedTime int
	foods := dinning_hall_elem.Foods

	for _, i := range currentOrder.Items {
		if foods[i-1].CookingApparatus == "" {
			A += foods[i-1].PreparationTime
		} else {
			C += foods[i-1].PreparationTime
		}
	}

	//Here i should use the variable from config and json config
	B = kitchenInform.CooksProfficiency
	D = kitchenInform.CookingApparatus
	E = kitchenInform.NrFoodsQueue
	F = len(currentOrder.Items)

	estimatedTime = (A/B + C/D) * (E + F) / F
	return estimatedTime

}

func SendOnlineOrder(receivedOrder OnlineReceivedOrder) KitchenInfo {
	reqBody, err := json.Marshal(receivedOrder)
	if err != nil {
		log.Fatal(err.Error())

	}
	//Sending the order to kitchen
	resp, err := http.Post(dinning_hall_elem.KitchenURL+"onlineOrder", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Request Failed: %s", err.Error())
	}
	defer resp.Body.Close()

	//receiving info from kitchen regarding foodQueue, apparatus, cooks profficiency,after the  online order is sent
	var kitchenInform KitchenInfo

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Response from Kitchen Failed: %s", err.Error())
		//return kitchenInform
	}
	_ = json.Unmarshal(body, &kitchenInform)

	log.Printf("---------The ONLINE order with id %d was SENT to Kitchen. Details: %+v\n", receivedOrder.Id, receivedOrder) // Unmarshal result
	return kitchenInform
}

func RegisterRestaurant() {
	var restaurInfo RestaurantInfo
	restaurInfo.RestaurantId = dinning_hall_elem.RestaurantId
	restaurInfo.Name = dinning_hall_elem.RestaurantName
	restaurInfo.Address = dinning_hall_elem.DinningHallUrl
	restaurInfo.MenuItems = dinning_hall_elem.NrFoods
	restaurInfo.Menu = dinning_hall_elem.Foods
	restaurInfo.Rating = dinning_hall_elem.Rating

	reqBody, err := json.Marshal(restaurInfo)
	if err != nil {
		log.Fatal(err.Error())

	}
	resp, err := http.Post(dinning_hall_elem.ManagerURL+"register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Registering Request Failed: %s", err.Error())
	}
	defer resp.Body.Close()
	log.Println("Successfully registered restaurant with id", restaurInfo.RestaurantId)
}
