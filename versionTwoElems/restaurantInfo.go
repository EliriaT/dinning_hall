package versionTwoElems

import dinning_hall_elem "github.com/EliriaT/dinning_hall/dinning-hall-elem"

type RestaurantInfo struct {
	RestaurantId int                      `json:"restaurant_id"`
	Name         string                   `json:"name"`
	Address      string                   `json:"address"`
	MenuItems    int                      `json:"menu_items"`
	Menu         []dinning_hall_elem.Food `json:"menu"`
	Rating       float32                  `json:"rating"`
}
