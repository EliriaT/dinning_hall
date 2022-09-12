package main

import (
	"math/rand"
	"sync"
)

const (
	maxFoods = 10
	minFoods = 1
	nrFoods  = 13
)

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

var aiOrder autoInc

type order struct {
	Id       int   `json:"id"`
	Items    []int `json:"items"`
	Priority int   `json:"priority"`
	MaxWait  int   `json:"max_wait"`
}

func newOrder() order {

	nrItems := rand.Intn(maxFoods) + minFoods
	foodList := make([]int, nrItems)
	maxWait := 0

	for i := 0; i < nrItems; i++ {
		foodId := rand.Intn(nrFoods) + 1
		//in the json the id starts from 1
		foodList[i] = foodId

		if prepTime := foods[foodId-1].PreparationTime; prepTime > maxWait {
			maxWait = prepTime
		}

	}
	return order{
		Id:       aiOrder.ID(),
		Items:    foodList,
		Priority: rand.Intn(5) + 1,
		MaxWait:  maxWait,
	}
}
