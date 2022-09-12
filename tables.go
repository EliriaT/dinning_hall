package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type tableState int

const (
	free tableState = iota
	waitToOrder
	waitToServe
)
const timeUnit = 50

// mutex to control when the table can generate order?
type table struct {
	id          int
	state       tableState
	clientOrder order
	lock        *sync.Mutex
}

var ordersChannel = make(chan int, len(tables))

func (t *table) makeOrder() {
	//t.lock.Lock() //unlock will be called when the order is served
	time.Sleep(timeUnit * time.Millisecond)
	t.state = waitToOrder
	t.clientOrder = newOrder()
	ordersChannel <- t.id //?
}

var tables = []table{
	{id: 1, state: free, lock: &sync.Mutex{}}, {id: 2, state: free, lock: &sync.Mutex{}},
	{id: 3, state: free, lock: &sync.Mutex{}}, {id: 4, state: free, lock: &sync.Mutex{}},
	{id: 5, state: free, lock: &sync.Mutex{}}, {id: 6, state: free, lock: &sync.Mutex{}},
	{id: 7, state: free, lock: &sync.Mutex{}}, {id: 8, state: free, lock: &sync.Mutex{}},
	{id: 9, state: free, lock: &sync.Mutex{}}, {id: 10, state: free, lock: &sync.Mutex{}},
}

func Init() {
	//defer wg.Done() IT CREATES A COPY OF TABLE IN T
	for i, _ := range tables {
		if rand.Intn(2) == 1 {
			tables[i].makeOrder()
			fmt.Printf("%+v of table %d \n", tables[i].clientOrder, tables[i].id)
		}
	}
	log.Printf("Init finished")
}
